package agent

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/do87/poly/src/logger"
	"github.com/do87/poly/src/polytree"
)

// agent is the agent service
type agent struct {
	MaxParallel  int              // Max plans running in parallel
	PlanTimeout  time.Duration    // Max time for running plans
	PollInterval time.Duration    // how often should the API be checked for new plan runs
	plans        map[string]*Plan // map of plan keys pointing to plans supported by the worker
	labels       Labels           // worker tags for filtering requests
	running      []string         // list of keys of running plans
	runLock      sync.Mutex
	hostname     string
}

// Plan is a type of a polytree
type Plan = polytree.Tree

// Job is a type of node in the polytree
type Job = polytree.Node

// Exec is the node execution function
type Exec = polytree.Exec

// Config represents the agent config
type Config struct {
	Labels Labels
	Key    Key
	Host   string
}

// Labels are the agent labels
type Labels map[string]string

// Key represents an agent key
type Key struct {
	Name       string
	PrivateKey []byte
}

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
		a.SetHost(c.Host)
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

// Register register plans to agent
func (a *agent) Register(plans ...*Plan) *agent {
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

	log.Info("agent is running...")

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
