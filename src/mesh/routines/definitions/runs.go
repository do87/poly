package definitions

import "context"

// AssignAgentToRuns handle agent assignment to created runs
func AssignAgentToRuns(ctx context.Context) {
}

// FindInactiveAgents checks for agents that didn't make liveness calls > 10 minutes
// and marks them as inactive
func FindInactiveAgents(ctx context.Context) {
}

// CancelRunsForInactiveAgents If an agent is marked as inactive but has a running job it needs to be marked as cancelled
func CancelRunsForInactiveAgents(ctx context.Context) {
}
