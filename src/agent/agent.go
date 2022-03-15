package agent

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/do87/poly/src/api/handlers/mesh/payloads"
	"github.com/do87/poly/src/api/handlers/mesh/present"
	"github.com/do87/poly/src/auth"
	"github.com/do87/poly/src/client"
	"github.com/do87/poly/src/logger"
	"github.com/do87/poly/src/polytree"
	"github.com/docker/distribution/uuid"
)

// agent is the agent service
type agent struct {
	MaxParallel     int              // Max plans running in parallel
	PlanTimeout     time.Duration    // Max time for running plans
	PollInterval    time.Duration    // how often should the API be checked for new plan runs
	uuid            uuid.UUID        // auto generated uuid for the agents
	plans           map[string]*Plan // map of plan keys pointing to plans supported by the worker
	labels          Labels           // worker tags for filtering requests
	running         []string         // list of keys of running plans
	runLock         sync.Mutex
	hostname        string
	client          *client.Client
	registrationKey auth.AgentRegisterKey
}

// Plan is a type of a polytree
type Plan = polytree.Tree

// Job is a type of node in the polytree
type Job = polytree.Node

// Exec is the node execution function
type Exec = polytree.Exec

// Config represents the agent config
type Config struct {
	Labels    Labels                // agent labels used to determine if the agent should run a pending plan execution
	Key       auth.AgentRegisterKey // key to register the agent with mesh API
	AgentHost string                // agent hostname
	MeshURL   string                // mesh API base URL
}

// Labels are the agent labels
type Labels []string

func (a *agent) execute(ctx context.Context, log *logger.Logger, plan *Plan, request *request) {
	t := (*polytree.Tree)(plan).Init()
	t.ExecuteWithTimeout(ctx, log, request.ID, request.Payload, a.PlanTimeout)
	a.done(ctx, log, plan)
}

func (a *agent) done(ctx context.Context, log *logger.Logger, plan *Plan) *agent { // TODO: tell API when execution finished
	log.Info("removing plan from agent running plans list", "plan", plan.Key)
	a.removeFromRunning(plan)
	return a
}

// removeFromRunning removes given plan from running lst
func (a *agent) removeFromRunning(plan *Plan) {
	newRunning := []string{}
	for _, run := range a.running {
		if run == plan.Key {
			continue
		}
		newRunning = append(newRunning, run)
	}

	a.runLock.Lock()
	defer a.runLock.Unlock()
	a.running = newRunning
}

// New returns a new agent
func New(cfg ...Config) *agent {
	a := &agent{
		MaxParallel:  3,
		PlanTimeout:  2 * time.Hour,
		PollInterval: 5 * time.Second,
		plans:        map[string]*Plan{},
		labels:       Labels{},
		hostname:     os.Getenv("HOST"),
	}
	for _, c := range cfg {
		a.SetLabels(c.Labels)
		a.SetHost(c.AgentHost)
		a.SetClient(c.MeshURL)
		a.SetKey(c.Key)
	}
	return a
}

// SetLabels sets agent labels
func (a *agent) SetLabels(l Labels) {
	for k, v := range l {
		a.labels[k] = v
	}
}

// SetHost sets agent host
func (a *agent) SetHost(host string) {
	if host != "" {
		a.hostname = host
	}
}

// SetClient sets client for mesh interactions
func (a *agent) SetClient(url string) {
	if url != "" {
		a.client = client.New(url)
	}
}

// SetKey sets agent registration key
func (a *agent) SetKey(key auth.AgentRegisterKey) {
	if key.PrivateKey != nil {
		a.registrationKey = key
	}
}

// Plans adds plans to agent
func (a *agent) Plans(plans ...*Plan) *agent {
	// add plans by their keys
	for _, plan := range plans {
		a.plans[plan.Key] = plan
	}
	return a
}

// Run runs the agent
func (a *agent) Run(ctx context.Context) {
	log, logsync := logger.New()
	defer logsync()

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	ticker := time.NewTicker(a.PollInterval)
	a.uuid = uuid.Generate()

	log.Info("registering agent...")
	res, err := a.registerAgent(ctx)
	if err != nil {
		panic(err)
	}
	log.Info("parsing answer...")
	if err := a.parseRegistrationResponse(res); err != nil {
		panic(err)
	}

	log.Info("agent is running.")
	for {
		select {
		case <-ticker.C:
			a.poll(ctx, log)
		case <-ctx.Done():
			a.eol(stop)
			return
		}
	}
}

func (a *agent) registerAgent(ctx context.Context) ([]byte, error) {
	if a.client == nil {
		return nil, errors.New("http client not configured correctly. make sure MeshURL is configured")
	}
	enc, err := a.registrationKey.Encode(a.hostname)
	if err != nil {
		return nil, err
	}
	return a.client.Do(ctx, http.MethodPost, "/agent", payloads.AgentRegister{
		UUID:     a.uuid.String(),
		Hostname: a.hostname,
		Labels:   a.labels,
		Plans:    a.getPlanNames(),
		EncodedKey: payloads.EncodedKey{
			Name:    a.registrationKey.Name,
			Encoded: enc,
		},
	})
}

func (a *agent) parseRegistrationResponse(r []byte) error {
	var res present.Presentor
	if err := json.Unmarshal(r, &res); err != nil {
		return err
	}
	v, ok := res.Data.(string)
	if !ok {
		return errors.New("failed to parse response data")
	}
	a.client.SetToken(v)
	return nil
}

func (a *agent) getPlanNames() []string {
	p := []string{}
	for _, v := range a.plans {
		p = append(p, v.Key)
	}
	return p
}

type request struct {
	ID      string
	Plan    string
	Payload []byte
}

// poll checks api for new plan requests
func (a *agent) poll(ctx context.Context, log *logger.Logger) {
	request := &request{
		ID:      "1",
		Plan:    "plan:infra:v1",
		Payload: []byte(`{"env": "dev"}`),
	}

	a.processRequest(ctx, log, request)
}

// processRequest find plan keys that match the request
func (a *agent) processRequest(ctx context.Context, log *logger.Logger, request *request) {
	var plan *Plan
	for _, p := range a.plans {
		if strings.EqualFold(p.Key, request.Plan) {
			plan = p
			break
		}
	}

	if plan == nil {
		return
	}

	if len(a.running) >= a.MaxParallel {
		return
	}

	a.runLock.Lock()
	defer a.runLock.Unlock()
	for _, key := range a.running {
		if key == plan.Key {
			return
		}
	}

	log.Info("Marking request as 'running'", "request", request)
	a.running = append(a.running, plan.Key)

	log.Info("Executing request", "request", request)
	go a.execute(ctx, log, plan, request)
}
