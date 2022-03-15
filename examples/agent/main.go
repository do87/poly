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
			"infra", "prod",
		},
		Key:       exampleKey,
		AgentHost: "localhost",
		MeshURL:   "127.0.0.1:8080",
	}).Plans(
		infra.Plan(),
	).Run(ctx)
}
