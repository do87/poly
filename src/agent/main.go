package main

import (
	"context"

	"github.com/do87/poly/src/agent/internal/agent"
)

func main() {
	agent.New(agent.Config{}).Run(context.Background())
}
