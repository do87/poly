package main

import (
	"context"

	"github.com/do87/poly/example/plans/infra"
	"github.com/do87/poly/src/agent"
)

func main() {
	ctx := context.Background()

	agent.New(agent.Labels{
		"infra": "prod",
	}).Register(
		infra.Plan(),
	).Run(ctx)
}
