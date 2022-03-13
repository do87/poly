package main

import (
	"github.com/do87/poly/src/api"
	"github.com/do87/poly/src/api/handlers/coordinator"
	"github.com/do87/poly/src/api/handlers/health"
)

func main() {
	api.New(api.Config{
		BindAddr: "127.0.0.1",
		DBConn:   "postgres://postgres:postgres@127.0.0.1:5432/poly?sslmode=disable",
	}).Register(
		health.Handler,
		coordinator.Handler,
	).Run()
}
