package agents

import (
	"github.com/do87/poly/src/db"
	"github.com/go-chi/chi/v5"
)

type Agents struct {
	route *chi.Mux
}

func Handler(r *chi.Mux, database *db.DB) {
	p := &Agents{
		route: r,
	}

	p.setRoutes()
}
