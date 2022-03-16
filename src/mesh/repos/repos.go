package repos

import (
	"github.com/do87/poly/src/db"
	"github.com/do87/poly/src/mesh/models"
	"gorm.io/gorm"
)

// Repo service
type Repo struct {
	Agents *agentsRepo
	Keys   *keysRepo
	Runs   *RunsRepo
}

type repo struct {
	db *gorm.DB
}

// New returns a new repo service
func New(d *db.DB) *Repo {
	if err := d.Migrate(
		models.Agent{},
		models.Key{},
		models.Run{},
	); err != nil {
		panic(err)
	}
	return &Repo{
		Agents: &agentsRepo{
			db: d.GetDB(),
		},
		Keys: &keysRepo{
			db: d.GetDB(),
		},
		Runs: &RunsRepo{
			db: d.GetDB(),
		},
	}
}
