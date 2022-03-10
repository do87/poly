package agents

// setRoutes attaches product routes
func (a *Agents) setRoutes() *Agents {
	a.route.Get("/agents", nil)           // list all agents
	a.route.Delete("/agent/{id}", nil)    // Delete an agent by ID
	a.route.Post("/agent/{id}/poll", nil) // agent API polling
	a.route.Get("/agents/keys", nil)      // list agent keys
	a.route.Post("/agents/key", nil)      // create a new agent key
	return a
}
