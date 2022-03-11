package usecases

import (
	"context"

	"github.com/do87/poly/src/api/handlers/agents/models"
)

// Usecase service
type Usecase struct {
	agents AgentsRepository
}

// AgentsRepository is the allowed usecase repo for agents
type AgentsRepository interface {
	List(ctx context.Context) ([]models.Agent, error)
}

// New creates a new Usecase service
func New(agents AgentsRepository) *Usecase {
	return &Usecase{
		agents: agents,
	}
}
