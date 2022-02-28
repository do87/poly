package agent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/do87/poly/sdk/job"
)

// execJob executes steps in a job
func (a *agent) execJob(ctx context.Context, result chan run, j *job.Job) error {
	defer func() { delete(a.plans, j.Key) }()

	var err error
	ch := make(chan error)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				ch <- fmt.Errorf("panic: %v", err)
			}
		}()
		for step := j.Get(); step != nil; {
			if err = ctx.Err(); err != nil {
				ch <- err
			}
			if step, err = step(); err != nil {
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
