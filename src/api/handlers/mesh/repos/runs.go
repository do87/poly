package repos

import (
	"context"

	"github.com/do87/poly/src/api/handlers/mesh/models"
	"gorm.io/gorm"
)

type RunsRepo struct {
	db *gorm.DB
}

// Get returns run by uuid
func (r *RunsRepo) Get(ctx context.Context, uuid string) (run models.Run, err error) {
	result := r.db.First(&run, "uuid = ?", uuid)
	if result.Error != nil {
		return models.Run{}, result.Error
	}
	return
}

// List returns all runs
func (r *RunsRepo) List(ctx context.Context) (keys []models.Run, err error) {
	result := r.db.Order("created_at DESC").Find(&keys)
	if result.Error != nil {
		return keys, result.Error
	}
	return keys, nil
}

// Create a new run
func (r *RunsRepo) Create(ctx context.Context, run models.Run) (models.Run, error) {
	if result := r.db.Create(&run); result.Error != nil {
		return run, result.Error
	}
	return run, nil
}

// Update updated a run by UUID
func (r *RunsRepo) Update(ctx context.Context, run models.Run) (models.Run, error) {
	m := &models.Run{}
	if err := r.db.Model(m).Where("uuid = ?", run.UUID).Updates(run).Error; err != nil {
		return run, err
	}
	return run, nil
}
