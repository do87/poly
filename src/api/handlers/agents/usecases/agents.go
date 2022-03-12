package usecases

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/do87/poly/src/api/handlers/agents/models"
	"github.com/do87/poly/src/api/handlers/agents/payloads"
)

// AgentsRepository is the allowed usecase repo for agents
type AgentsRepository interface {
	List(ctx context.Context) ([]models.Agent, error)
	Register(ctx context.Context, agent models.Agent) (models.Agent, error)
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
func (u *agentsUsecase) Register(ctx context.Context, r *http.Request) (agent models.Agent, err error) {
	var payload payloads.Agent
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return
	}
	return u.repo.Register(ctx, payload.ToModel())
}
