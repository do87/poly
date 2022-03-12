package health

import (
	"github.com/do87/poly/src/db"
	"github.com/go-chi/chi/v5"
)

type health struct {
	route *chi.Mux
}

func Handler(r *chi.Mux, database *db.DB) {
	p := &health{
		route: r,
	}

	p.setRoutes()
}
