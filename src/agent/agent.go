package agent

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/do87/poly/src/mesh/api/payloads"
	"github.com/do87/poly/src/mesh/api/present"
	"github.com/do87/poly/src/mesh/common"
	"github.com/do87/poly/src/pkg/auth"
	"github.com/do87/poly/src/pkg/client"
	"github.com/do87/poly/src/pkg/logger"
	"github.com/do87/poly/src/pkg/polytree"
	"github.com/docker/distribution/uuid"
)

// agent is the agent service
type agent struct {
	MaxParallel     int              // Max plans running in parallel
	PlanTimeout     time.Duration    // Max time for running plans
	PollInterval    time.Duration    // how often should the API be checked for new plan runs
	uuid            uuid.UUID        // auto generated uuid for the agents
	log             logger.Log       // logger
	plans           map[string]*Plan // map of plan keys pointing to plans supported by the worker
	labels          Labels           // worker tags for filtering requests
	running         []string         // list of keys of running requests
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
	Logger    logger.Log            // the logger
}

// Labels are the agent labels
type Labels []string

// New returns a new agent
func New(c Config) *agent {
	validateConfig(c)
	a := &agent{
		MaxParallel:     3,
		PlanTimeout:     2 * time.Hour,
		PollInterval:    5 * time.Second,
		plans:           map[string]*Plan{},
		labels:          c.Labels,
		hostname:        os.Getenv("HOST"),
		client:          client.New(c.MeshURL),
		registrationKey: c.Key,
		log:             c.Logger,
	}

	// override host if specified
	a.SetHost(c.AgentHost)
	return a
}

func validateConfig(c Config) {
	if c.Logger == nil {
		panic("logger not configured")
	}
	if c.Key.PrivateKey == nil {
		panic("registration key not configured")
	}
}

// SetHost sets agent host
func (a *agent) SetHost(host string) {
	if host != "" {
		a.hostname = host
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

func (a *agent) getPlanKeys() []string {
	keys := []string{}
	for _, plan := range a.plans {
		keys = append(keys, plan.Key)
	}
	return keys
}

func (a *agent) logAgentInfo() {
	if a.hostname != "" {
		a.log.Info("- Hostname: " + a.hostname)
	}
	a.log.Info("- Agent UUID: " + a.uuid.String())
	a.log.Info("- Labels: " + strings.Join(a.labels, ", "))
	a.log.Info("- Supported Plans: " + strings.Join(a.getPlanKeys(), ", "))
}

// Run runs the agent
func (a *agent) Run(ctx context.Context) {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	ticker := time.NewTicker(a.PollInterval)
	a.uuid = uuid.Generate()
	a.log.Info("⏫ agent starting up...")
	a.logAgentInfo()

	a.log.Info("🔐 registering agent...")
	res, err := a.registerAgent(ctx)
	if err != nil {
		panic(err)
	}
	if err := a.parseRegistrationResponse(res); err != nil {
		panic(err)
	}

	a.log.Info("🚀 agent is running")
	for {
		select {
		case <-ticker.C:
			a.poll(ctx, a.log)
		case <-ctx.Done():
			a.eol(a.log, stop)
			return
		}
	}
}

func (a *agent) registerAgent(ctx context.Context) ([]byte, error) {
	if a.client == nil {
		return nil, errors.New("http client not configured correctly. make sure MeshURL is configured")
	}
	enc, err := a.registrationKey.Encode(a.uuid.String())
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
	v, err := present.GetDataFromBytes[string](r)
	if err != nil {
		return err
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
func (a *agent) poll(ctx context.Context, log logger.Log) {
	log.Debug("checking for pending runs...")
	b, err := a.client.Do(ctx, http.MethodGet, fmt.Sprintf("/agent/%s/runs/pending", a.uuid), nil)
	if err != nil {
		log.Error(err.Error())
		return
	}

	runs, err := present.GetDataFromBytes[[]present.TypeRun](b)
	if err != nil {
		log.Error(err.Error())
		return
	}
	for _, run := range runs {
		request := &request{
			ID:      run.UUID,
			Plan:    run.Plan,
			Payload: []byte(run.Payload),
		}
		a.processRequest(ctx, log, request)
	}
}

// processRequest find plan keys that match the request
func (a *agent) processRequest(ctx context.Context, log logger.Log, request *request) {
	log.Info(fmt.Sprintf("📋 processing run id %s with plan %s", request.ID, request.Plan), "payload", request.Payload)
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

	a.markRequestAsRunning(ctx, log, request)

	log.Info(fmt.Sprintf("🧑‍🔧 executing run id %s", request.ID))
	go a.execute(ctx, log, request, plan)
}

func (a *agent) markRequestAsRunning(ctx context.Context, log logger.Log, request *request) {
	log.Info(fmt.Sprintf("⏱  marking run id %s as running", request.ID))

	// api call
	if err := a.setRunStatus(ctx, request, common.RUN_STATUS_RUNNING); err != nil {
		log.Error(err.Error())
		return
	}

	// append to running list
	a.running = append(a.running, request.ID)
}

func (a *agent) setRunStatus(ctx context.Context, request *request, status string) (err error) {
	_, err = a.client.Do(ctx,
		http.MethodPatch,
		fmt.Sprintf("/run/%s", request.ID),
		payloads.RunUpdate{Status: status},
	)
	return
}

func (a *agent) execute(ctx context.Context, log logger.Log, request *request, plan *Plan) {
	t := (*polytree.Tree)(plan).Init()
	t.ExecuteWithTimeout(ctx, log, request.ID, request.Payload, a.PlanTimeout)
	a.done(ctx, log, request, t.HasErrors())
}

func (a *agent) done(ctx context.Context, log logger.Log, request *request, hasErrors bool) *agent { // TODO: tell API when execution finished
	text := "successfully"
	status := common.RUN_STATUS_SUCCESS
	if hasErrors {
		text = "with errors"
		status = common.RUN_STATUS_ERROR
	}
	log.Info(fmt.Sprintf("🚦 run id %s finished %s", request.ID, text))
	log.Info(fmt.Sprintf("📤 finalizing run id %s", request.ID))
	a.removeFromRunning(request)
	if err := a.setRunStatus(ctx, request, status); err != nil {
		log.Error(err.Error())
		return a
	}
	log.Info("✅ done", "request", request)
	return a
}

// removeFromRunning removes given request from running lst
func (a *agent) removeFromRunning(request *request) {
	newRunning := []string{}
	for _, run := range a.running {
		if run == request.ID {
			continue
		}
		newRunning = append(newRunning, run)
	}
	a.runLock.Lock()
	defer a.runLock.Unlock()
	a.running = newRunning
}
