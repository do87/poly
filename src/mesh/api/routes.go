package mesh

import (
	"github.com/do87/poly/src/mesh/api/middlewares"
	"github.com/go-chi/chi/v5"
)

// setRoutes attaches product routes
func (h *handler) setRoutes() *handler {
	// agents lifecycle
	h.route.Get("/agents", h.agents.list(h.agents.uc))                // list all agents
	h.route.Post("/agent", h.agents.register(h.agents.uc, h.keys.uc)) // agent registration

	// keys
	h.route.Get("/agents/keys", h.keys.list(h.keys.uc))            // list agent keys
	h.route.Post("/agents/key", h.keys.create(h.keys.uc))          // create a new agent key
	h.route.Delete("/agents/key/{name}", h.keys.delete(h.keys.uc)) // Delete an agent key by ID

	h.route.Group(func(r chi.Router) {
		r.Use(middlewares.VerifyAgent)

		// agents lifecycle
		h.route.Get("/agent/{id}/runs/pending", h.runs.listPending(h.runs.uc, h.agents.uc)) // List pending runs
		h.route.Delete("/agent/{id}", h.agents.deregister(h.agents.uc))                     // Deregisters an agent by ID

		// runs lifecycle
		h.route.Get("/runs", h.runs.list(h.runs.uc))
		h.route.Post("/run", h.runs.create(h.runs.uc))
		h.route.Put("/run/{id}", h.runs.create(h.runs.uc))
	})

	return h
}
