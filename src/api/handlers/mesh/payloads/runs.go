package payloads

import "github.com/do87/poly/src/api/handlers/mesh/models"

// RunCreate is the payload needed to create an agent run
type RunCreate struct {
	UUID   string   `json:"id"`
	Labels []string `json:"labels"`
	Plan   string   `json:"plan"`
}

// ToModel converts payload to model
func (a *RunCreate) ToModel() models.Run {
	return models.Run{
		UUID:   a.UUID,
		Labels: a.Labels,
		Plan:   a.Plan,
	}
}
