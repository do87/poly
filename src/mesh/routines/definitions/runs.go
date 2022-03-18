package definitions

import (
	"context"

	"github.com/do87/poly/src/mesh/common"
	"github.com/do87/poly/src/mesh/models"
	"github.com/do87/poly/src/mesh/usecases"
	"github.com/do87/poly/src/pkg/logger"
)

// AssignAgentToRuns handle agent assignment to created runs
func AssignAgentToRuns(ctx context.Context, log logger.Log, r, a *usecases.Usecase) {
	for _, err := range r.Runs.AssignAgentToRuns(ctx, a.Agents, assignmentProcess) {
		log.Error(err.Error())
	}
}

func assignmentProcess(runs []models.Run, agents []models.Agent) map[string]string {
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

func filterCandidates(candidates map[string][]string) map[string]string {
	counters := map[string]int{}
	flattened := map[string]string{}
	for runID, agents := range candidates {
		flattened[runID] = ""
		for _, a := range agents {
			counters[a]++
		}
		if len(agents) == 0 {
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
		flattened[runID] = min
	}
	return flattened
}

// CancelRunsForInactiveAgents If an agent is marked as inactive but has a running job it needs to be marked as cancelled
func CancelRunsForInactiveAgents(ctx context.Context, log logger.Log, r, a *usecases.Usecase) {
	if err := r.Runs.CancelRunsForInactiveAgents(ctx, a.Agents); err != nil {
		log.Error(err.Error())
	}
}
