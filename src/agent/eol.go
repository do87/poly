package agent

import (
	"context"
	"net/http"
	"os"

	"github.com/do87/poly/src/pkg/logger"
)

// eol handles agent end of life
func (a *agent) eol(log logger.Log, stop context.CancelFunc) {
	log.Info("agent end of life:")
	// code to run before agent stops
	a.destructor(log)

	stop()
	os.Exit(0)
}

// destructor contains logic that runs right before the agent is stopped
func (a *agent) destructor(log logger.Log) {
	log.Info("- marking agent as inactive")
	if _, err := a.client.Do(context.TODO(), http.MethodDelete, "/agent/"+a.uuid.String(), nil); err != nil {
		log.Error(err.Error())
	}
}
