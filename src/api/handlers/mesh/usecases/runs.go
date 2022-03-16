package usecases

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/do87/poly/src/api/handlers/mesh/common"
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

// Create registers an key and returns it
func (u *runsUsecase) Create(ctx context.Context, r *http.Request) (key models.Run, err error) {
	var payload payloads.RunCreate
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return
	}
	model := payload.ToModel()
	if err = common.SetRunStatus(&model, common.RUN_STATUS_PENDING); err != nil {
		return
	}
	return u.repo.Create(ctx, model)
}
