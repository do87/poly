package usecases

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/do87/poly/src/mesh/api/payloads"
	"github.com/do87/poly/src/mesh/common"
	"github.com/do87/poly/src/mesh/models"
)

// AgentsRepository is the allowed usecase repo for agents
type AgentsRepository interface {
	Get(ctx context.Context, id string) (agents models.Agent, err error)
	List(ctx context.Context) ([]models.Agent, error)
	ListActive(ctx context.Context) ([]models.Agent, error)
	ListInactiveSince(ctx context.Context, t time.Time) ([]models.Agent, error)
	Register(ctx context.Context, agent models.Agent) (models.Agent, error)
	Deregister(ctx context.Context, id string) (models.Agent, error)
	Update(ctx context.Context, agent models.Agent) (models.Agent, error)
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

// ListActive returns a list of all active agents
func (u *agentsUsecase) ListActive(ctx context.Context, r *http.Request) ([]models.Agent, error) {
	return u.repo.ListActive(ctx)
}

// Register registers an agent and returns it
func (u *agentsUsecase) Register(ctx context.Context, r *http.Request, k *keysUsecase) (agent models.Agent, err error) {
	var payload payloads.AgentRegister
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return
	}
	key, err := k.repo.Get(ctx, payload.EncodedKey.Name)
	if err != nil {
		return
	}
	if err = common.ProcessRegisterKey(key, payload.EncodedKey.Encoded, payload.UUID); err != nil {
		return
	}
	return u.repo.Register(ctx, payload.ToModel(key.Name))
}

// Deregister unregisters an agent
func (u *agentsUsecase) Deregister(ctx context.Context, r *http.Request, id string) (models.Agent, error) {
	return u.repo.Deregister(ctx, id)
}

// Ping updates the agent's updated_at field
func (u *agentsUsecase) Ping(ctx context.Context, r *http.Request, id string) (agent models.Agent, err error) {
	agent, err = u.repo.Get(ctx, id)
	if err != nil {
		return
	}
	agent.UpdatedAt = time.Now()
	return u.repo.Update(ctx, agent)
}

// MarkInactiveAgents checks for agents that didn't make liveness calls > 10 minutes
// and marks them as inactive
func (u *agentsUsecase) MarkInactiveAgents(ctx context.Context) error {
	agents, err := u.repo.ListInactiveSince(ctx, time.Now().Add(-time.Minute*10))
	if err != nil {
		return err
	}
	for _, a := range agents {
		a.Active = false
		if _, err := u.repo.Update(ctx, a); err != nil {
			return err
		}
	}
	return nil
}
