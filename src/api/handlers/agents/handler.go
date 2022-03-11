package agents

import (
	"github.com/do87/poly/src/api/handlers/agents/repos"
	"github.com/do87/poly/src/db"
	"github.com/go-chi/chi/v5"
)

type Agents struct {
	route *chi.Mux
	repo  *repos.Repo
}

func Handler(r *chi.Mux, d *db.DB) {
	p := &Agents{
		route: r,
		repo:  repos.New(d),
	}
	p.setRoutes()
}
