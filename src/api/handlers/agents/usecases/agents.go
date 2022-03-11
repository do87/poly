package usecases

import (
	"context"
	"net/http"

	"github.com/do87/poly/src/api/handlers/agents/models"
)

type agentsUsecase struct {
	repo AgentsRepository
}

// List returns a list of all agents
func (u *agentsUsecase) List(ctx context.Context, r *http.Request) ([]models.Agent, error) {
	return u.repo.List(ctx)
}
