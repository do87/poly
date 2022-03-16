package main

import (
	"github.com/do87/poly/src/api"
	"github.com/do87/poly/src/health"
	meshAPI "github.com/do87/poly/src/mesh/api"
)

func main() {
	api.New(api.Config{
		BindAddr: "127.0.0.1",
		DBConn:   "postgres://postgres:postgres@127.0.0.1:5432/poly?sslmode=disable",
	}).Register(
		health.Handler,
		meshAPI.Handler,
	).Run()
}
