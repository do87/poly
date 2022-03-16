package routines

import (
	"time"

	"github.com/do87/poly/src/mesh/routines/definitions"
	"github.com/go-co-op/gocron"
)

// Run runs all configured routines
func Run() {
	cron := gocron.NewScheduler(time.UTC)
	cron.Every(2).Second().Do(definitions.AssignAgentToRuns)
	cron.Every(30).Second().Do(definitions.FindInactiveAgents)
	cron.Every(1).Minute().Do(definitions.CancelRunsForInactiveAgents)
	cron.StartBlocking()
}
