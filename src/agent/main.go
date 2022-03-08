package main

import (
	"context"

	"github.com/do87/poly/src/agent/internal/agent"
	"github.com/do87/poly/src/agent/internal/logger"
)

func main() {
	ctx := context.Background()
	log, logsync := logger.New()
	defer logsync()

	agent.Register(agent.Labels{
		"infra": "prod",
	}).Run(ctx, log)
}
