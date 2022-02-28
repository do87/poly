package agent

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/do87/poly/src/agent/internal/polytree"
)

// agent is the agent service
type agent struct {
	MaxParallel  uint8         // Max plans running in parallel
	PlanTimeout  time.Duration // Max time for running plans
	PollInterval time.Duration // how often should the API be checked for new plan runs
	plans        map[string]*Plan
	tags         Tags
}

type Plan polytree.Tree
type Tags map[string]string

// Register creates a new worker and sets its tags and plans
func Register(tags Tags, plans ...*Plan) *agent {
	a := &agent{
		MaxParallel:  3,
		PlanTimeout:  2 * time.Hour,
		PollInterval: 5 * time.Second,
		plans:        map[string]*Plan{},
		tags:         tags,
	}

	// add plans by their keys
	for _, plan := range plans {
		a.plans[plan.Key] = plan
	}
	return a
}

// Run runs the agent
func (a *agent) Run(ctx context.Context) {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	ticker := time.NewTicker(a.PollInterval)

	select {
	case <-ticker.C:
		a.poll(ctx)
		return
	case <-ctx.Done():
		a.eol(stop)
		return
	}
}

// poll checks api for new plan requests
func (a *agent) poll(ctx context.Context) {
	request := "plan1" // TODO
	if plan := a.processRequest(request); plan != nil {
		pt := (*polytree.Tree)(plan)
		pt.Execute(ctx) // TODO: check parallelism, Add timeout
	}
}

func (a *agent) processRequest(planRequest string) *Plan {
	for _, plan := range a.plans {
		if plan.Key == planRequest {
			return plan
		}
	}
	return nil
}
