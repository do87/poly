package definitions

import (
	"context"
	"time"

	"github.com/do87/poly/src/mesh/repos"
	"github.com/do87/poly/src/pkg/logger"
)

// FindInactiveAgents checks for agents that didn't make liveness calls > 10 minutes
// and marks them as inactive
func FindInactiveAgents(ctx context.Context, log *logger.Logger, r *repos.Repo) {
	agents, err := r.Agents.ListAgentsSinceFilterByActive(ctx, true, time.Now().Add(time.Minute*time.Duration(-10)))
	if err != nil {
		log.Error(err.Error())
		return
	}
	for _, a := range agents {
		a.Active = false
		if _, err := r.Agents.Update(ctx, a); err != nil {
			log.Error(err.Error())
		}
	}
}
