package routines

import (
	"context"
	"time"

	"github.com/do87/poly/src/mesh/repos"
	"github.com/do87/poly/src/mesh/routines/definitions"
	"github.com/do87/poly/src/mesh/usecases"
	"github.com/do87/poly/src/pkg/db"
	"github.com/do87/poly/src/pkg/logger"
	"github.com/go-co-op/gocron"
)

// Run runs all configured routines
func Run(ctx context.Context, log logger.Log, db *db.DB) {
	log.Info("♻️  starting mesh server routines")
	cron := gocron.NewScheduler(time.UTC)

	r := repos.New(db)
	agents := usecases.NewAgentsUsecase(r.Agents)
	cron.Every(2).Second().SingletonMode().Do(definitions.AssignAgentToRuns, ctx, log, agents)
	cron.Every(30).Second().Do(definitions.FindInactiveAgents, ctx, log, r)
	cron.Every(1).Minute().Do(definitions.CancelRunsForInactiveAgents, ctx, log, r)
	cron.StartBlocking()
}
