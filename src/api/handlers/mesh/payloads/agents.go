package payloads

import "github.com/do87/poly/src/mesh/models"

// AgentRegister is the payload needed to register an agent
type AgentRegister struct {
	UUID       string     `json:"id"`
	Hostname   string     `json:"hostname"`
	Labels     []string   `json:"labels"`
	Plans      []string   `json:"plans"`
	EncodedKey EncodedKey `json:"encoded_key"`
}

// ToModel converts payload to model
func (a *AgentRegister) ToModel(name string) models.Agent {
	return models.Agent{
		UUID:     a.UUID,
		KeyName:  name,
		Hostname: a.Hostname,
		Labels:   a.Labels,
		Plans:    a.Plans,
	}
}
