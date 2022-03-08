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
	p := agent.NewPlan("plan:infra:v1", &infra{}, 1*time.Hour)

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
	log.Info("Handling Terraform State Storage")
	return nil, nil
}

func tfRunNode(ctx context.Context, log *logger.Logger, meta interface{}, payload []byte) (polytree.Exec, error) {
	log.Info("Handling Terraform Run")
	return nil, nil
}
