package agent

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// agent is the agent service
type agent struct {
	config Config
}

// Config holds the agent's configuration
type Config struct {
	MaxParallel uint8         // Max plans running in parallel
	PlanTimeout time.Duration // Max time for running plans
	Ticker      time.Duration // how often should the API be checked for new plan runs
}

// New returns a new agent service
func New(c Config) *agent {
	a := &agent{}
	return a.setConfig(c)
}

// setConfig sets the agent configuration with the given config
// and sets default values if not specified
func (a *agent) setConfig(c Config) *agent {
	if c.Ticker == 0 {
		c.Ticker = 5 * time.Second
	}
	if c.MaxParallel == 0 {
		c.MaxParallel = 3
	}
	if c.PlanTimeout == 0 {
		c.Ticker = 2 * time.Hour
	}
	a.config = c
	return a
}

// Run runs the agent
func (a *agent) Run(ctx context.Context) {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	ticker := time.NewTicker(a.config.Ticker)

	select {
	case <-ticker.C:
		a.findPlanRequests(ctx)
		return
	case <-ctx.Done():
		a.eol(stop)
		return
	}
}

// findPlanRequests checks api for new plan requests
func (a *agent) findPlanRequests(ctx context.Context) {
	{

	}
}

// handlePostrun handles operation following a plan run
func (a *agent) postrun(ctx context.Context, planKey string) {

	// run failed

}
