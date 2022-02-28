package agent

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/do87/poly/sdk/job"
	"github.com/do87/poly/sdk/plan"
)

// agent is the agent service
type agent struct {
	config Config
	plans  map[string]*plan.Plan
}

type run struct {
	job *job.Job
	err error
}

// Config holds the agent's configuration
type Config struct {
	Ticker      time.Duration
	MaxParallel uint8
	PlanTimeout time.Duration
}

// New returns a new agent service
func New(c Config) *agent {
	a := &agent{}
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
func (a *agent) Run(ctx context.Context) {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	ticker := time.NewTicker(a.config.Ticker)

	select {
	case <-ticker.C:
		if len(a.plans) < int(a.config.MaxParallel) {
			a.findPlanRequests(ctx)
		}
		return
	case <-ctx.Done():
		a.eol(stop)
		return
	}
}

// findPlanRequests checks api for new plan requests
func (a *agent) findPlanRequests(ctx context.Context) {
	{
		p := plan.New()
		a.plans[p.Key] = p
		a.runJobs(ctx, p.Jobs)

	}
}

func (a *agent) runJobs(ctx context.Context, jobs []*job.Job) error {
	numJobs := len(jobs)
	var results chan run = make(chan run, numJobs)

	for _, job := range jobs {
		go a.execJob(ctx, results, job)
	}

	for i := 0; i < numJobs; i++ {
		res := <-results
		if res.err != nil {
			return res.err
		}

		if len(res.job.Children) > 0 {
			return a.runJobs(ctx, res.job.Children)
		}
	}

	return nil
}
