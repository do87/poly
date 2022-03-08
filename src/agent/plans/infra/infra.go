package infra

import (
	"context"
	"time"

	"github.com/do87/poly/src/agent/internal/logger"
	"github.com/do87/poly/src/agent/internal/polytree"
)

type infra struct{}

func Plan() *polytree.Tree {
	p := &polytree.Tree{
		Key:     "plan:infra:v1",
		Meta:    &infra{},
		Timeout: 1 * time.Hour,
	}

	node := &polytree.Node{
		Key:  "state-storage",
		Exec: polytree.Exec(stateStorageNode),
	}

	p.AddNode(node)
	return p
}

func stateStorageNode(ctx context.Context, log *logger.Logger, meta interface{}, payload []byte) (polytree.Exec, error) {
	log.Info("[Infra Plan] State Storage Node")
	return nil, nil
}
