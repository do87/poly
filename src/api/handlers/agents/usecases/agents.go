package usecases

import (
	"context"
	"net/http"

	"github.com/do87/poly/src/api/handlers/agents/models"
)

// AgentsRepository is the allowed usecase repo for agents
type AgentsRepository interface {
	List(ctx context.Context) ([]models.Agent, error)
	Register(ctx context.Context) (models.Agent, error)
}

// NewAgentsUsecase creates a new Usecase service
func NewAgentsUsecase(agents AgentsRepository) *Usecase {
	return &Usecase{
		Agents: &agentsUsecase{
			repo: agents,
		},
	}
}

type agentsUsecase struct {
	repo AgentsRepository
}

// List returns a list of all agents
func (u *agentsUsecase) List(ctx context.Context, r *http.Request) ([]models.Agent, error) {
	return u.repo.List(ctx)
}

// Register registers an agent and returns it
func (u *agentsUsecase) Register(ctx context.Context, r *http.Request) (models.Agent, error) {
	return u.repo.Register(ctx)
}
