package usecases

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/do87/poly/src/api/handlers/mesh/models"
	"github.com/do87/poly/src/api/handlers/mesh/payloads"
)

// AgentsRepository is the allowed usecase repo for agents
type AgentsRepository interface {
	List(ctx context.Context) ([]models.Agent, error)
	Register(ctx context.Context, agent models.Agent) (models.Agent, error)
	Deregister(ctx context.Context, id string) error
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
func (u *agentsUsecase) Register(ctx context.Context, r *http.Request, keysUc *keysUsecase) (agent models.Agent, err error) {
	var payload payloads.AgentRegister
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return
	}
	key, err := keysUc.repo.GetByName(ctx, payload.EncodedKey.Name)
	if err != nil {
		return
	}
	if err = u.processKey(key, payload); err != nil {
		return
	}
	return u.repo.Register(ctx, payload.ToModel(key.UUID))
}

// Deregister unregisters an agent
func (u *agentsUsecase) Deregister(ctx context.Context, r *http.Request, id string) (err error) {
	return u.repo.Deregister(ctx, id)
}
