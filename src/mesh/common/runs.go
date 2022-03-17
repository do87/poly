package common

import (
	"errors"
	"time"

	"github.com/do87/poly/src/mesh/models"
	"github.com/docker/distribution/uuid"
)

// Run status related costs
const (
	RUN_STATUS_CREATED   = "created"
	RUN_STATUS_PENDING   = "pending"
	RUN_STATUS_RUNNING   = "running"
	RUN_STATUS_SUCCESS   = "success"
	RUN_STATUS_ERROR     = "error"
	RUN_STATUS_CANCELED  = "canceled"
	RUN_STATUS_NO_AGENTS = "no_agents"
)

// SetRunStatus modifies a given model according to the provided status
func SetRunStatus(run *models.Run, status string) error {
	if err := validateRunStatus(status); err != nil {
		return err
	}
	run.Status = status
	switch status {
	case RUN_STATUS_CREATED:
		handleRunStatusCreated(run)
	case RUN_STATUS_PENDING:
		handleRunStatusPending(run)
	case RUN_STATUS_RUNNING:
		handleRunStatusRunning(run)
	case RUN_STATUS_NO_AGENTS:
		fallthrough
	case RUN_STATUS_CANCELED:
		fallthrough
	case RUN_STATUS_ERROR:
		fallthrough
	case RUN_STATUS_SUCCESS:
		handleRunStatusDone(run)
	}
	return nil
}

func handleRunStatusCreated(run *models.Run) {
	run.UUID = uuid.Generate().String()
	run.CreatedAt = time.Now()
}

func handleRunStatusPending(run *models.Run) {
	run.AssignedAt = time.Now()
}

func handleRunStatusRunning(run *models.Run) {
	run.StartedAt = time.Now()
}

func handleRunStatusDone(run *models.Run) {
	run.FinishedAt = time.Now()
}

func validateRunStatus(status string) error {
	switch status {
	case RUN_STATUS_CANCELED:
		fallthrough
	case RUN_STATUS_CREATED:
		fallthrough
	case RUN_STATUS_ERROR:
		fallthrough
	case RUN_STATUS_PENDING:
		fallthrough
	case RUN_STATUS_RUNNING:
		fallthrough
	case RUN_STATUS_SUCCESS:
		fallthrough
	case RUN_STATUS_NO_AGENTS:
		return nil
	}
	return errors.New("invalid run status")
}
