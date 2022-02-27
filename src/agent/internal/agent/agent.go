package agent

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/do87/poly/sdk/plan"
)

// agent is the agent service
type agent struct {
	config Config
	plans  map[string]*plan.Plan
	run    chan string
}

// Config holds the agent's configuration
type Config struct {
	Ticker      time.Duration
	MaxParallel uint8
	PlanTimeout time.Duration
}

// New returns a new agent service
func New(c Config) *agent {
	a := &agent{
		run: make(chan string),
	}
	return a.setConfig(c)
}

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
	return a
}

// Run runs the agent
func (a *agent) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	ticker := time.NewTicker(a.config.Ticker)

	select {
	case <-ticker.C:
		if len(a.plans) < int(a.config.MaxParallel) {
			a.findPlanRequests()
		}
		return
	case <-ctx.Done():
		a.eol(stop)
		return
	}
}

// findPlanRequests checks api for new plan requests
func (a *agent) findPlanRequests() {
	{
		p := plan.New()
		a.plans[p.Key] = p
		a.run <- p.Key
		go a.execute(p)
	}
}
