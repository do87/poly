package main

import (
	"context"

	"github.com/do87/poly/src/agent/internal/agent"
	"github.com/do87/poly/src/agent/plans/infra"
)

func main() {
	ctx := context.Background()

	agent.New(agent.Labels{
		"infra": "prod",
	}).Register(
		infra.Plan(),
	).Run(ctx)
}
