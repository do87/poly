package main

import (
	"context"
	"os"

	"github.com/do87/poly/example/agent/plans/infra"
	"github.com/do87/poly/src/agent"
	"github.com/do87/poly/src/pkg/auth"
	"github.com/do87/poly/src/pkg/logger"
)

func main() {
	ctx := context.Background()
	log, logsync := logger.NewDevelopment()
	defer logsync()

	agent.New(agent.Config{
		Labels: agent.Labels{
			"infra", "prod",
		},
		Key: auth.AgentRegisterKey{
			Name:       "pubkey:v0.1",
			PrivateKey: getPrivateKey(),
		},
		AgentHost: "localhost",
		MeshURL:   "http://127.0.0.1:8080",
	}).Plans(
		infra.Plan(),
	).Run(ctx, log)
}

func getPrivateKey() []byte {
	k, err := os.ReadFile("keys/private_key.pem")
	if err != nil {
		panic(err)
	}
	return k
}
