package mesh

import (
	"github.com/do87/poly/src/api/handlers/mesh/repos"
	"github.com/do87/poly/src/api/handlers/mesh/usecases"
	"github.com/do87/poly/src/db"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	route  *chi.Mux
	repo   *repos.Repo
	agents *agents
	keys   *keys
}

// Handler handles agent related routes and functionality
func Handler(r *chi.Mux, d *db.DB) {
	repo := repos.New(d)
	p := &handler{
		route: r,
		repo:  repo,

		// Link usecases:
		agents: &agents{
			uc: usecases.NewAgentsUsecase(repo.Agents),
		},
		keys: &keys{
			uc: usecases.NewKeysUsecase(repo.Keys),
		},
	}
	p.setRoutes()
}
