package definitions

import (
	"context"

	"github.com/do87/poly/src/mesh/usecases"
	"github.com/do87/poly/src/pkg/logger"
)

// MarkInactiveAgents checks for agents that didn't make liveness calls > 10 minutes
// and marks them as inactive
func MarkInactiveAgents(ctx context.Context, log logger.Log, u *usecases.Usecase) {
	if err := u.Agents.MarkInactiveAgents(ctx); err != nil {
		log.Error(err.Error())
	}
}
