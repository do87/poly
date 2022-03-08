package health

import (
	"net/http"

	"github.com/go-chi/render"
)

// setRoutes attaches product routes
func (p *Health) setRoutes() *Health {
	p.route.Get("/.well-known/live", live())
	return p
}

func live() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]string{"status": "ok"})
	}
}
