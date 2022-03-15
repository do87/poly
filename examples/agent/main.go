package main

import (
	"context"
	"os"

	"github.com/do87/poly/example/agent/plans/infra"
	"github.com/do87/poly/src/agent"
	"github.com/do87/poly/src/auth"
)

func main() {
	ctx := context.Background()

	k, err := os.ReadFile("private_key.pem")
	if err != nil {
		panic(err)
	}

	agent.New(agent.Config{
		Labels: agent.Labels{
			"infra", "prod",
		},
		Key: auth.AgentRegisterKey{
			Name:       "pubkey:v0.1",
			PrivateKey: k,
		},
		AgentHost: "localhost",
		MeshURL:   "http://127.0.0.1:8080",
	}).Plans(
		infra.Plan(),
	).Run(ctx)
}
