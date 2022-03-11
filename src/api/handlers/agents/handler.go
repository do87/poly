package agents

import (
	"github.com/do87/poly/src/api/handlers/agents/repos"
	"github.com/do87/poly/src/api/handlers/agents/usecases"
	"github.com/do87/poly/src/db"
	"github.com/go-chi/chi/v5"
)

type Agents struct {
	route *chi.Mux
	repo  *repos.Repo
	uc    *usecases.Usecase
}

func Handler(r *chi.Mux, d *db.DB) {
	repo := repos.New(d)
	p := &Agents{
		route: r,
		repo:  repo,
		uc:    usecases.New(repo.Agents),
	}
	p.setRoutes()
}
