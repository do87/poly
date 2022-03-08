package agent

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/do87/poly/src/agent/internal/logger"
	"github.com/do87/poly/src/agent/internal/polytree"
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

type Plan polytree.Tree
type Labels map[string]string

func (a *agent) execute(ctx context.Context, log *logger.Logger, plan *Plan, payload []byte) {
	t := (*polytree.Tree)(plan)
	t.Payload = payload
	t.Timeout = a.PlanTimeout
	t.ExecuteWithTimeout(ctx, log, a.PlanTimeout)
	a.save(ctx, plan)
}

func (a *agent) save(ctx context.Context, plan *Plan) *agent { // TODO
	return a
}

// Register creates a new worker and sets its tags and plans
func Register(labels Labels, plans ...*Plan) *agent {
	a := &agent{
		MaxParallel:  3,
		PlanTimeout:  2 * time.Hour,
		PollInterval: 5 * time.Second,
		plans:        map[string]*Plan{},
		labels:       labels,
	}

	// add plans by their keys
	for _, plan := range plans {
		a.plans[plan.Key] = plan
	}
	return a
}

// Run runs the agent
func (a *agent) Run(ctx context.Context, log *logger.Logger) {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	ticker := time.NewTicker(a.PollInterval)

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

// poll checks api for new plan requests
func (a *agent) poll(ctx context.Context, log *logger.Logger) {
	request := "plan:infra:v1" // TODO
	payload := []byte(`{"env": "dev"}`)
	a.processRequest(ctx, log, request, payload)
}

// processRequest find plan keys that match the request
func (a *agent) processRequest(ctx context.Context, log *logger.Logger, planRequest string, payload []byte) {
	var plan *Plan
	for _, p := range a.plans {
		if strings.EqualFold(p.Key, planRequest) {
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
	a.running = append(a.running, plan.Key)
	a.runLock.Unlock()

	a.execute(ctx, log, plan, payload)
}
