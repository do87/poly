package agent

import (
	"context"
	"os"
	"time"
)

// eol handles agent end of life
func (a *agent) eol(stop context.CancelFunc) {
	// code to run before agent stops
	close(a.run)
	a.destructor()

	stop()
	os.Exit(0)
}

// destructor contains logic that runs right before the agent is stopped
func (a *agent) destructor() {
	for len(a.plans) > 0 {
		time.Sleep(time.Second * 30)
	}
}
