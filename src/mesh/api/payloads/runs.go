package payloads

import (
	"github.com/do87/poly/src/mesh/common"
	"github.com/do87/poly/src/mesh/models"
	"github.com/docker/distribution/uuid"
)

// RunCreate is the payload needed to create an agent run
type RunCreate struct {
	Labels  []string `json:"labels"`
	Plan    string   `json:"plan"`
	Payload string   `json:"payload"`
}

// ToModel converts payload to model
func (a *RunCreate) ToModel() (m models.Run, err error) {
	m = models.Run{
		UUID:    uuid.Generate().String(),
		Labels:  a.Labels,
		Plan:    a.Plan,
		Payload: []byte(a.Payload),
	}
	err = common.SetRunStatus(&m, common.RUN_STATUS_CREATED)
	return
}

// RunUpdate is the payload needed to update an agent run
type RunUpdate struct {
	Status string `json:"status"`
}

// ToModel converts payload to model
func (a *RunUpdate) ToModel(uuid string) (m models.Run, err error) {
	m = models.Run{
		UUID: uuid,
	}
	err = common.SetRunStatus(&m, a.Status)
	return
}
