package health

import (
	"github.com/go-chi/chi/v5"
)

type Health struct {
	route *chi.Mux
}

func Handle(r *chi.Mux) {
	p := &Health{
		route: r,
	}

	p.setRoutes()
}
