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

// RunsRepository is the usecase repo for agent runs
type RunsRepository interface {
	Get(ctx context.Context, name string) (models.Run, error)
	List(ctx context.Context) ([]models.Run, error)
	ListCreated(ctx context.Context) ([]models.Run, error)
	ListPendingByAgentID(ctx context.Context, agentID string) ([]models.Run, error)
	ListPendingSince(ctx context.Context, t time.Time) (keys []models.Run, err error)
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

// List returns a list of all runs
func (u *runsUsecase) ListPendingByAgentID(ctx context.Context, r *http.Request, agentID string) ([]models.Run, error) {
	return u.repo.ListPendingByAgentID(ctx, agentID)
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
func (u *runsUsecase) Update(ctx context.Context, r *http.Request, uuid string) (run models.Run, err error) {
	var payload payloads.RunUpdate
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return
	}
	if run, err = payload.ToModel(uuid); err != nil {
		return
	}
	return u.repo.Update(ctx, run)
}

// AssignAgentToRuns handle agent assignment to created runs
// assignFn decides which agents is assigned to execute a run and returns
// map[string][]string, where the key is the run UUID and the value is the agent UUID or empty
// if no agent is found for executing the run
func (u *runsUsecase) AssignAgentToRuns(ctx context.Context, a *agentsUsecase, assignFn func([]models.Run, []models.Agent) map[string]string) []error {
	runs, err := u.repo.ListCreated(ctx)
	if err != nil {
		return []error{err}
	}
	if len(runs) == 0 {
		return nil
	}

	// get active agents
	agents, err := a.repo.ListActive(ctx)
	if err != nil {
		return []error{err}
	}

	candidates := assignFn(runs, agents)
	selected := []string{}
	el := []error{}
	for _, run := range runs {
		if candidates[run.UUID] == "" {
			if err = u.updateRunStatus(ctx, run, common.RUN_STATUS_NO_AGENTS); err != nil {
				el = append(el, err)
			}
			continue
		}
		run.Agent = candidates[run.UUID]
		selected = append(selected, run.Agent)
		if err := u.updateRunStatus(ctx, run, common.RUN_STATUS_PENDING); err != nil {
			el = append(el, err)
		}
	}

	for _, agent := range agents {
		if !common.Contains(selected, agent.UUID) {
			continue
		}
		agent.LastAssignedAt = time.Now()
		if _, err = a.repo.Update(ctx, agent); err != nil {
			el = append(el, err)
		}
	}
	if len(el) == 0 {
		return nil
	}
	return el
}

func (u *runsUsecase) updateRunStatus(ctx context.Context, run models.Run, status string) error {
	if err := common.SetRunStatus(&run, status); err != nil {
		return err
	}
	if _, err := u.repo.Update(ctx, run); err != nil {
		return err
	}
	return nil
}

// CancelRunsForInactiveAgents If an agent is marked as inactive but has a running job it needs to be marked as cancelled
func (u *runsUsecase) CancelRunsForInactiveAgents(ctx context.Context, a *agentsUsecase) error {
	runs, err := u.repo.ListPendingSince(ctx, time.Now().Add(-time.Minute*10))
	if err != nil {
		return err
	}

	runsToCancel := u.findRunsToCancel(ctx, runs, a)
	for _, run := range runsToCancel {
		if err := u.updateRunStatus(ctx, run, common.RUN_STATUS_CANCELED); err != nil {
			return err
		}
	}
	return nil
}

func (u *runsUsecase) findRunsToCancel(ctx context.Context, runs []models.Run, a *agentsUsecase) []models.Run {
	toCancel := []models.Run{}
	for _, run := range runs {
		if run.Agent == "" {
			toCancel = append(toCancel, run)
			continue
		}
		a, err := a.repo.Get(ctx, run.Agent)
		if err != nil {
			continue
		}
		if !a.Active {
			toCancel = append(toCancel, run)
		}
	}
	return toCancel
}
