package definitions

import (
	"context"
	"time"

	"github.com/do87/poly/src/mesh/common"
	"github.com/do87/poly/src/mesh/models"
	"github.com/do87/poly/src/mesh/repos"
	"github.com/do87/poly/src/pkg/logger"
)

// AssignAgentToRuns handle agent assignment to created runs
func AssignAgentToRuns(ctx context.Context, log *logger.Logger, r *repos.Repo) {
	runs, err := r.Runs.ListCreated(ctx)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if len(runs) == 0 {
		return
	}

	// get active agents
	agents, err := r.Agents.ListActive(ctx)
	if err != nil {
		log.Error(err.Error())
		return
	}

	candidates := findAgentsForRuns(runs, agents)
	selected := []string{}
	for _, run := range runs {
		if len(candidates[run.UUID]) == 0 {
			if err := updateRun(ctx, r, run, common.RUN_STATUS_NO_AGENTS); err != nil {
				log.Error(err.Error())
			}
			continue
		}
		run.Agent = candidates[run.UUID][0]
		selected = append(selected, run.Agent)
		if err := updateRun(ctx, r, run, common.RUN_STATUS_PENDING); err != nil {
			log.Error(err.Error())
		}
	}

	for _, agent := range agents {
		if !common.Contains(selected, agent.UUID) {
			continue
		}
		agent.LastAssignedAt = time.Now()
		if _, err := r.Agents.Update(ctx, agent); err != nil {
			log.Error(err.Error())
		}
	}

}

func updateRun(ctx context.Context, r *repos.Repo, run models.Run, status string) error {
	if err := common.SetRunStatus(&run, status); err != nil {
		return err
	}
	if _, err := r.Runs.Update(ctx, run); err != nil {
		return err
	}
	return nil
}

func findAgentsForRuns(runs []models.Run, agents []models.Agent) map[string][]string {
	candidates := map[string][]string{}
	for _, run := range runs {
		candidates[run.UUID] = []string{}
		for _, agent := range agents {
			if !common.Contains(agent.Plans, run.Plan) {
				continue
			}
			if !common.SubsetOf(run.Labels, agent.Labels) {
				continue
			}
			candidates[run.UUID] = append(candidates[run.UUID], agent.UUID)
		}
	}
	return filterCandidates(candidates)
}

func filterCandidates(candidates map[string][]string) map[string][]string {
	counters := map[string]int{}
	for runID, agents := range candidates {
		for _, a := range agents {
			counters[a]++
		}
		if len(agents) <= 1 {
			continue
		}
		min := agents[0]
		for _, a := range agents {
			if counters[a] < counters[min] {
				min = a
			}
		}
		for _, a := range agents {
			if a != min {
				counters[a]--
			}
		}
		candidates[runID] = []string{min}
	}
	return candidates
}

// CancelRunsForInactiveAgents If an agent is marked as inactive but has a running job it needs to be marked as cancelled
func CancelRunsForInactiveAgents(ctx context.Context, log *logger.Logger, r *repos.Repo) {
	runs, err := r.Runs.ListPendingSince(ctx, time.Now().Add(time.Minute*time.Duration(-10)))
	if err != nil {
		log.Error(err.Error())
		return
	}

	runsToCancel := findRunsToCancel(ctx, log, r, runs)
	for _, run := range runsToCancel {
		if err := updateRun(ctx, r, run, common.RUN_STATUS_CANCELED); err != nil {
			log.Error(err.Error())
		}
	}
}

func findRunsToCancel(ctx context.Context, log *logger.Logger, r *repos.Repo, runs []models.Run) []models.Run {
	toCancel := []models.Run{}
	for _, run := range runs {
		if run.Agent == "" {
			log.Warning("found a run in pending state with no agent uuid", "run", run.UUID)
			toCancel = append(toCancel, run)
			continue
		}
		a, err := r.Agents.Get(ctx, run.Agent)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		if !a.Active {
			toCancel = append(toCancel, run)
		}
	}
	return toCancel
}
