package infra

import (
	"context"
	"time"

	"github.com/do87/poly/src/agent"
	"github.com/do87/poly/src/logger"
	"github.com/do87/poly/src/polytree"
)

type infra struct{}

func Plan() *agent.Plan {
	p := &agent.Plan{
		Key:     "plan:infra:v1",
		Meta:    &infra{},
		Timeout: 1 * time.Hour,
	}

	storage := &polytree.Node{
		Key:  "state-storage",
		Exec: polytree.Exec(stateStorageNode),
	}

	tfrun := &polytree.Node{
		Key:  "run-terraform",
		Exec: polytree.Exec(tfRunNode),
	}

	p.AddNode(storage)
	p.AddNode(tfrun)
	p.ParentOf(storage, tfrun)
	return p
}

func stateStorageNode(ctx context.Context, log *logger.Logger, meta interface{}, payload []byte) (polytree.Exec, error) {
	log.Info("[Infra Plan] State Storage Node")
	return nil, nil
}

func tfRunNode(ctx context.Context, log *logger.Logger, meta interface{}, payload []byte) (polytree.Exec, error) {
	log.Info("[Infra Plan] Run Terraform Node")
	return nil, nil
}
