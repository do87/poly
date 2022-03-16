package infra

import (
	"context"
	"time"

	"github.com/do87/poly/src/agent"
	"github.com/do87/poly/src/pkg/logger"
)

type infra struct{}

func Plan() *agent.Plan {
	p := &agent.Plan{
		Key:     "plan:infra:v1",
		Meta:    &infra{},
		Timeout: 1 * time.Hour,
	}

	storage := &agent.Job{
		Key:  "state-storage",
		Exec: agent.Exec(stateStorageNode),
	}

	tfrun := &agent.Job{
		Key:  "run-terraform",
		Exec: agent.Exec(tfRunNode),
	}

	p.AddNode(storage)
	p.AddNode(tfrun)
	p.ParentOf(storage, tfrun)
	return p
}

func stateStorageNode(ctx context.Context, log *logger.Logger, meta interface{}, payload []byte) (agent.Exec, error) {
	log.Info("Handling Terraform State Storage")
	return nil, nil
}

func tfRunNode(ctx context.Context, log *logger.Logger, meta interface{}, payload []byte) (agent.Exec, error) {
	log.Info("Handling Terraform Run")
	return nil, nil
}
