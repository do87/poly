package agents

// setRoutes attaches product routes
func (a *handler) setRoutes() *handler {
	a.route.Get("/agents", a.agents.list(a.uc))              // list all agents
	a.route.Post("/agent", a.agents.register(a.uc))          // agent registration
	a.route.Delete("/agent/{id}", a.agents.deregister(a.uc)) // Deregisters an agent by ID

	a.route.Get("/agents/keys", nil)        // list agent keys
	a.route.Post("/agents/key", nil)        // create a new agent key
	a.route.Delete("/agents/key/{id}", nil) // Delete an agent key by ID
	return a
}
