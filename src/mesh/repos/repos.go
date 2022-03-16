package repos

import (
	"github.com/do87/poly/src/mesh/models"
	"github.com/do87/poly/src/pkg/db"
	"gorm.io/gorm"
)

// Repo service
type Repo struct {
	Agents *agentsRepo
	Keys   *keysRepo
	Runs   *runsRepo
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
		Runs: &runsRepo{
			db: d.GetDB(),
		},
	}
}
