package present

import (
	"fmt"
	"time"

	"github.com/amalfra/etag"
	"github.com/do87/poly/src/api/handlers/agents/models"
)

const (
	KEY_AGENT  = "poly:agent"
	KEY_AGENTS = "poly:agents"
)

type agent struct {
	UUID      string    `json:"id"`
	Hostname  string    `json:"hostname"`
	Labels    []string  `json:"labels"`
	Plans     []string  `json:"plans"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (agent) FromModel(m models.Agent) agent {
	return agent{
		UUID:      m.UUID,
		Hostname:  m.Hostname,
		Labels:    m.Labels,
		Plans:     m.Plans,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func Agents(agentModels []models.Agent) Presentor {
	p := make([]agent, 0)
	a := agent{}
	u := ""
	for _, m := range agentModels {
		u += fmt.Sprintf("%s-%s-%s;", m.UUID, m.Hostname, m.UpdatedAt.String())
		p = append(p, a.FromModel(m))
	}
	return wrap(KEY_AGENTS, etag.Generate(u, true), p)
}

func Agent(m models.Agent) Presentor {
	a := agent{}
	u := fmt.Sprintf("%s-%s-%s;", m.UUID, m.Hostname, m.UpdatedAt.String())
	return wrap(KEY_AGENTS, etag.Generate(u, true), a.FromModel(m))
}
