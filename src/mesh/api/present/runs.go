package present

import (
	"fmt"
	"time"

	"github.com/amalfra/etag"
	"github.com/do87/poly/src/mesh/models"
)

// key constants
const (
	KEY_RUN  = "poly:run"
	KEY_RUNS = "poly:runs"
)

type run struct {
	UUID       string    `json:"id"`
	Agent      string    `json:"agent_id"`
	Plan       string    `json:"plan"`
	Labels     []string  `json:"labels"`
	Status     string    `json:"status"`
	Payload    string    `json:"payload"`
	AssignedAt time.Time `json:"assigned_at"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// FromModel converts model to presenter struct
func (run) FromModel(m models.Run) run {
	return run{
		UUID:       m.UUID,
		Agent:      m.Agent,
		Plan:       m.Plan,
		Labels:     m.Labels,
		Payload:    string(m.Payload),
		Status:     m.Status,
		AssignedAt: m.AssignedAt,
		StartedAt:  m.StartedAt,
		FinishedAt: m.FinishedAt,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

var _r = run{}

// Runs presents agent runs
func Runs(models []models.Run) Presentor {
	p := make([]run, 0)
	u := ""
	for _, m := range models {
		u += fmt.Sprintf("%s-%s;", m.UUID, m.UpdatedAt.String())
		p = append(p, _r.FromModel(m))
	}
	return wrap(KEY_RUNS, etag.Generate(u, true), p)
}

// Run presents an agent run
func Run(m models.Run) Presentor {
	u := fmt.Sprintf("%s-%s;", m.UUID, m.UpdatedAt.String())
	return wrap(KEY_RUN, etag.Generate(u, true), _r.FromModel(m))
}
