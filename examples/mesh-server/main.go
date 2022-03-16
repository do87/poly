package main

import (
	"github.com/do87/poly/src/mesh"
)

func main() {
	mesh.New(mesh.Config{
		API: mesh.APIConfig{
			BindAddr: "127.0.0.1",
		},
		DBConn: "postgres://postgres:postgres@127.0.0.1:5432/poly?sslmode=disable",
	}).Run()
}
