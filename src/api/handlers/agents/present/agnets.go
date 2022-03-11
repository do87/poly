package present

import (
	"fmt"
	"time"

	"github.com/amalfra/etag"
	"github.com/do87/poly/src/api/handlers/agents/models"
)

const (
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

func Agents(agents []models.Agent) Presentor {
	p := make([]agent, 0)
	u := ""
	for _, a := range agents {
		u += fmt.Sprintf("%s-%s-%s;", a.UUID, a.Hostname, a.UpdatedAt.String())
		p = append(p, agent{
			UUID:      a.UUID,
			Hostname:  a.Hostname,
			Labels:    a.Labels,
			Plans:     a.Plans,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		})
	}
	return wrap(KEY_AGENTS, etag.Generate(u, true), p)
}
