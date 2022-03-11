package agents

import (
	"net/http"

	"github.com/do87/poly/src/api/handlers/agents/usecases"
	"github.com/go-chi/render"
)

// setRoutes attaches product routes
func (a *Agents) setRoutes() *Agents {
	a.route.Get("/agents", listAgents(a.uc)) // list all agents
	a.route.Post("/agent/{id}/poll", nil)    // agent API polling
	a.route.Delete("/agent/{id}", nil)       // Delete an agent by ID

	a.route.Get("/agents/keys", nil)        // list agent keys
	a.route.Post("/agents/key", nil)        // create a new agent key
	a.route.Delete("/agents/key/{id}", nil) // Delete an agent key by ID
	return a
}

func listAgents(u *usecases.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := u.Agents.List(r.Context(), r)
		if err != nil {
			render.JSON(w, r, err)
		}
		render.JSON(w, r, data)
	}
}
