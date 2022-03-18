package present

import (
	"fmt"
	"time"

	"github.com/amalfra/etag"
	"github.com/do87/poly/src/mesh/models"
)

// key consts
const (
	KEY_AGENT              = "poly:agent"
	KEY_AGENT_ACCESS_TOKEN = "poly:agent:access-token"
	KEY_AGENTS             = "poly:agents"
)

type agent struct {
	UUID      string    `json:"id"`
	Hostname  string    `json:"hostname"`
	Active    bool      `json:"active"`
	Labels    []string  `json:"labels"`
	Plans     []string  `json:"plans"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FromModel converts model to presenter struct
func (agent) FromModel(m models.Agent) agent {
	return agent{
		UUID:      m.UUID,
		Hostname:  m.Hostname,
		Active:    m.Active,
		Labels:    m.Labels,
		Plans:     m.Plans,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

var _a = agent{}

// Agents presents agents
func Agents(agentModels []models.Agent) Presentor {
	p := make([]agent, 0)
	u := ""
	for _, m := range agentModels {
		u += fmt.Sprintf("%s-%s-%s;", m.UUID, m.Hostname, m.UpdatedAt.String())
		p = append(p, _a.FromModel(m))
	}
	return wrap(KEY_AGENTS, etag.Generate(u, true), p)
}

// Agent presents agent
func Agent(m models.Agent) Presentor {
	u := fmt.Sprintf("%s-%s-%s;", m.UUID, m.Hostname, m.UpdatedAt.String())
	return wrap(KEY_AGENTS, etag.Generate(u, true), _a.FromModel(m))
}

// AccessToken presents agent's access token
func AccessToken(t string) Presentor {
	return wrap(KEY_AGENT_ACCESS_TOKEN, etag.Generate(t, true), t)
}
