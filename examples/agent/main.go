package main

import (
	"context"

	"github.com/do87/poly/example/agent/plans/infra"
	"github.com/do87/poly/src/agent"
)

func main() {
	ctx := context.Background()

	agent.New(agent.Config{
		Labels: agent.Labels{
			"infra": "prod",
		},
		Key: exampleKey,
	}).Register(
		infra.Plan(),
	).Run(ctx)
}
