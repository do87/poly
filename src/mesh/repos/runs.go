package repos

import (
	"context"
	"time"

	"github.com/do87/poly/src/mesh/common"
	"github.com/do87/poly/src/mesh/models"
)

// runsRepo is the repo for runs
type runsRepo repo

// Get returns run by uuid
func (r *runsRepo) Get(ctx context.Context, uuid string) (run models.Run, err error) {
	result := r.db.First(&run, "uuid = ?", uuid)
	if result.Error != nil {
		return models.Run{}, result.Error
	}
	return
}

// List returns all runs
func (r *runsRepo) List(ctx context.Context) (keys []models.Run, err error) {
	result := r.db.Order("created_at DESC").Find(&keys)
	if result.Error != nil {
		return keys, result.Error
	}
	return keys, nil
}

// List pending runs
func (r *runsRepo) ListPending(ctx context.Context, agentID string) (keys []models.Run, err error) {
	result := r.db.Where("status = ?", common.RUN_STATUS_PENDING).Where("agent = ?", agentID).Order("created_at DESC").Find(&keys)
	if result.Error != nil {
		return keys, result.Error
	}
	return keys, nil
}

// ListCreated returns all new runs with unassigned agent
func (r *runsRepo) ListCreated(ctx context.Context) (keys []models.Run, err error) {
	result := r.db.Where("status = ?", common.RUN_STATUS_CREATED).Order("created_at DESC").Find(&keys)
	if result.Error != nil {
		return keys, result.Error
	}
	return keys, nil
}

// ListPendingSince returns all runs in pending state since given time
func (r *runsRepo) ListPendingSince(ctx context.Context, t time.Time) (keys []models.Run, err error) {
	result := r.db.
		Where("status = ?", common.RUN_STATUS_PENDING).
		Where("assigned_at > ?", t).
		Order("created_at DESC").Find(&keys)
	if result.Error != nil {
		return keys, result.Error
	}
	return keys, nil
}

// Create a new run
func (r *runsRepo) Create(ctx context.Context, run models.Run) (models.Run, error) {
	if result := r.db.Create(&run); result.Error != nil {
		return run, result.Error
	}
	return run, nil
}

// Update updated a run by UUID
func (r *runsRepo) Update(ctx context.Context, run models.Run) (models.Run, error) {
	m := &models.Run{}
	if err := r.db.Model(m).Where("uuid = ?", run.UUID).Updates(run).Error; err != nil {
		return run, err
	}
	return run, nil
}
