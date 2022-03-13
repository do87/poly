package repos

import (
	"context"

	"github.com/do87/poly/src/api/handlers/mesh/models"
	"gorm.io/gorm"
)

type agentsRepo struct {
	db *gorm.DB
}

// Get returns agent by UUID
func (r *agentsRepo) Get(ctx context.Context, id string) (agents models.Agent, err error) {
	var agent models.Agent
	result := r.db.First(&agent, "uuid = ?", id)
	if result.Error != nil {
		return models.Agent{}, result.Error
	}
	return agents, nil
}

// List returns all agents
func (r *agentsRepo) List(ctx context.Context) (agents []models.Agent, err error) {
	result := r.db.Order("uuid ASC").Find(&agents)
	if result.Error != nil {
		return agents, result.Error
	}
	return agents, nil
}

// Register registers the agent
func (r *agentsRepo) Register(ctx context.Context, agent models.Agent) (models.Agent, error) {
	if result := r.db.FirstOrCreate(&agent); result.Error != nil {
		return agent, result.Error
	}
	return agent, nil
}

// Deregister deletes the agent by uuid
func (r *agentsRepo) Deregister(ctx context.Context, id string) error {
	a, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := r.db.Delete(&a).Error; err != nil {
		return err
	}
	return nil
}
