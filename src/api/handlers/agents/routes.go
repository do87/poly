package agents

// setRoutes attaches product routes
func (a *handler) setRoutes() *handler {
	a.route.Get("/agents", a.agents.list(a.agents.uc))              // list all agents
	a.route.Post("/agent", a.agents.register(a.agents.uc))          // agent registration
	a.route.Delete("/agent/{id}", a.agents.deregister(a.agents.uc)) // Deregisters an agent by ID

	a.route.Get("/agents/keys", a.keys.list(a.keys.uc))            // list agent keys
	a.route.Post("/agents/key", a.keys.create(a.agents.uc))        // create a new agent key
	a.route.Delete("/agents/key/{id}", a.keys.delete(a.agents.uc)) // Delete an agent key by ID
	return a
}
