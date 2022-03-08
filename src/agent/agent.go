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
}

type Plan = polytree.Tree

type Labels map[string]string

func (a *agent) execute(ctx context.Context, log *logger.Logger, plan *Plan, request *Request) {
	t := (*polytree.Tree)(plan)
	t = t.Copy()
	t.ExecuteWithTimeout(ctx, request.ID, request.Payload, a.PlanTimeout)
	a.done(ctx, plan)
}

func (a *agent) done(ctx context.Context, plan *Plan) *agent { // TODO: tell API when execution finished
	return a
}

// New returens a new agent
func New(labels Labels) *agent {
	a := &agent{
		MaxParallel:  3,
		PlanTimeout:  2 * time.Hour,
		PollInterval: 5 * time.Second,
		plans:        map[string]*Plan{},
		labels:       labels,
	}
	return a
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

type Request struct {
	ID      string
	Plan    string
	Payload []byte
}

// poll checks api for new plan requests
func (a *agent) poll(ctx context.Context, log *logger.Logger) {
	request := &Request{
		ID:      "1",
		Plan:    "plan:infra:v1",
		Payload: []byte(`{"env": "dev"}`),
	}

	a.processRequest(ctx, log, request)
}

// processRequest find plan keys that match the request
func (a *agent) processRequest(ctx context.Context, log *logger.Logger, request *Request) {
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

	log.Info("Marking request as 'running'", "request", request)
	a.runLock.Lock()
	a.running = append(a.running, plan.Key)
	a.runLock.Unlock()

	log.Info("Executing request", "request", request)
	a.execute(ctx, log, plan, request)
}
