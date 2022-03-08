package main

import (
	"github.com/do87/poly/src/api"
	"github.com/do87/poly/src/api/handlers/health"
)

func main() {
	api.New(api.Config{}).Register(
		health.Handler,
	).Run()
}
