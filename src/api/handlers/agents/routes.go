package agents

// setRoutes attaches product routes
func (a *handler) setRoutes() *handler {
	a.route.Get("/agents", a.agents.list(a.uc))         // list all agents
	a.route.Post("/agent", a.agents.register(a.uc))     // agent registration
	a.route.Post("/agent/{id}/poll", nil)               // agent API polling
	a.route.Delete("/agent", a.agents.deregister(a.uc)) // Deregisters an agent by UUID & Hostname

	a.route.Get("/agents/keys", nil)        // list agent keys
	a.route.Post("/agents/key", nil)        // create a new agent key
	a.route.Delete("/agents/key/{id}", nil) // Delete an agent key by ID
	return a
}
