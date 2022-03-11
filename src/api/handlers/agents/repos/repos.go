package repos

import (
	"github.com/do87/poly/src/api/handlers/agents/models"
	"github.com/do87/poly/src/db"
)

// Repo service
type Repo struct {
	Agents *agentsRepo
}

// New returns a new repo service
func New(d *db.DB) *Repo {
	if err := d.Migrate(
		models.Agent{},
	); err != nil {
		panic(err)
	}
	return &Repo{
		Agents: &agentsRepo{
			db: d.GetDB(),
		},
	}
}
