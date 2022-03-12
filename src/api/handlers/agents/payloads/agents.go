package payloads

import "github.com/do87/poly/src/api/handlers/agents/models"

// AgentRegister is the payload needed to register an agent
type AgentRegister struct {
	UUID     string   `json:"id"`
	Hostname string   `json:"hostname"`
	Labels   []string `json:"labels"`
	Plans    []string `json:"plans"`
}

// ToModel converts payload to model
func (a *AgentRegister) ToModel() models.Agent {
	return models.Agent{
		UUID:     a.UUID,
		Hostname: a.Hostname,
		Labels:   a.Labels,
		Plans:    a.Plans,
	}
}
