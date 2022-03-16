package usecases

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/do87/poly/src/api/handlers/mesh/models"
	"github.com/do87/poly/src/api/handlers/mesh/payloads"
)

// RunsRepository is the usecase repo for agent runs
type RunsRepository interface {
	Get(ctx context.Context, name string) (models.Run, error)
	List(ctx context.Context) ([]models.Run, error)
	Create(ctx context.Context, run models.Run) (models.Run, error)
	Update(ctx context.Context, run models.Run) (models.Run, error)
}

// NewRunsUsecase creates a new Usecase service
func NewRunsUsecase(runs RunsRepository) *Usecase {
	return &Usecase{
		Runs: &runsUsecase{
			repo: runs,
		},
	}
}

type runsUsecase struct {
	repo RunsRepository
}

// List returns a list of all runs
func (u *runsUsecase) List(ctx context.Context, r *http.Request) ([]models.Run, error) {
	return u.repo.List(ctx)
}

// Create creates a new run
func (u *runsUsecase) Create(ctx context.Context, r *http.Request) (run models.Run, err error) {
	var payload payloads.RunCreate
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return
	}
	if run, err = payload.ToModel(); err != nil {
		return
	}
	return u.repo.Create(ctx, run)
}

// Update updates an existing run
func (u *runsUsecase) Update(ctx context.Context, r *http.Request) (run models.Run, err error) {
	var payload payloads.RunCreate
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return
	}
	if run, err = payload.ToModel(); err != nil {
		return
	}
	return u.repo.Update(ctx, run)
}
