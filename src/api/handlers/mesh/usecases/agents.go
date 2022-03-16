package usecases

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/do87/poly/src/api/handlers/mesh/payloads"
	"github.com/do87/poly/src/mesh/common"
	"github.com/do87/poly/src/mesh/models"
)

// AgentsRepository is the allowed usecase repo for agents
type AgentsRepository interface {
	List(ctx context.Context) ([]models.Agent, error)
	Register(ctx context.Context, agent models.Agent) (models.Agent, error)
	Deregister(ctx context.Context, id string) (models.Agent, error)
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
	key, err := keysUc.repo.Get(ctx, payload.EncodedKey.Name)
	if err != nil {
		return
	}
	if err = common.ProcessRegisterKey(key, payload.EncodedKey.Encoded, payload.Hostname); err != nil {
		return
	}
	return u.repo.Register(ctx, payload.ToModel(key.Name))
}

// Deregister unregisters an agent
func (u *agentsUsecase) Deregister(ctx context.Context, r *http.Request, id string) (models.Agent, error) {
	return u.repo.Deregister(ctx, id)
}
