package main

import (
	"context"

	"github.com/do87/poly/src/agent/internal/agent"
)

func main() {
	ctx := context.Background()
	agent.Register(agent.Tags{
		"infra": "prod",
	}).Run(ctx)
}
