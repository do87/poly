package main

import (
	"context"

	"github.com/do87/poly/src/mesh"
	"github.com/do87/poly/src/pkg/logger"
)

func main() {
	ctx := context.Background()
	log, logsync := logger.NewDevelopment()
	defer logsync()

	mesh.New(mesh.Config{
		Logger: log,
		API: mesh.APIConfig{
			BindAddr: "127.0.0.1",
			BindPort: 8080,
		},
		DBConn: "postgres://postgres:postgres@127.0.0.1:5432/poly?sslmode=disable",
	}).Run(ctx)
}
