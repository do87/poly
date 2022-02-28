package agent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/do87/poly/sdk/plan"
)

// execute initiates a plan run
func (a *agent) execute(ctx context.Context, p *plan.Plan) error {
	<-a.run
	defer func() { delete(a.plans, p.Key) }()

	var err error
	ch := make(chan error)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				ch <- fmt.Errorf("panic: %v", err)
			}
		}()
		for job := p.Start(); job != nil; {
			if err = ctx.Err(); err != nil {
				ch <- err
			}
			if job, err = job(); err != nil {
				ch <- err
			}
		}
		ch <- nil
	}()

	select {
	case err := <-ch:
		return err
	case <-time.After(a.config.PlanTimeout):
		return errors.New("node execution timed out")
	}

}
