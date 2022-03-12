package repos

import (
	"context"

	"github.com/do87/poly/src/api/handlers/agents/models"
	"gorm.io/gorm"
)

type agentsRepo struct {
	db *gorm.DB
}

// List returns all products
func (r *agentsRepo) List(ctx context.Context) (agents []models.Agent, err error) {
	result := r.db.Order("uuid ASC").Find(&agents)
	if result.Error != nil {
		return agents, result.Error
	}
	return agents, nil
}

// Register registers the agent
func (r *agentsRepo) Register(ctx context.Context) (agent models.Agent, err error) {
	return agent, nil
}
