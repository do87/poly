package agent

import (
	"context"
	"os"
)

// eol handles agent end of life
func (a *agent) eol(stop context.CancelFunc) {
	// code to run before agent stops
	a.destructor()

	stop()
	os.Exit(0)
}

// destructor contains logic that runs right before the agent is stopped
func (a *agent) destructor() {

}
