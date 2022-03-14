package mesh

import (
	"github.com/do87/poly/src/api/handlers/mesh/middlewares"
	"github.com/go-chi/chi/v5"
)

// setRoutes attaches product routes
func (h *handler) setRoutes() *handler {
	h.route.Get("/agents", h.agents.list(h.agents.uc))                // list all agents
	h.route.Post("/agent", h.agents.register(h.agents.uc, h.keys.uc)) // agent registration

	h.route.Group(func(r chi.Router) {
		r.Use(middlewares.VerifyAgent)
		h.route.Delete("/agent/{id}", h.agents.deregister(h.agents.uc)) // Deregisters an agent by ID
	})

	h.route.Get("/agents/keys", h.keys.list(h.keys.uc))          // list agent keys
	h.route.Post("/agents/key", h.keys.create(h.keys.uc))        // create a new agent key
	h.route.Delete("/agents/key/{id}", h.keys.delete(h.keys.uc)) // Delete an agent key by ID
	return h
}
