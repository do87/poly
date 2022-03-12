package agents

import (
	"github.com/do87/poly/src/api/handlers/agents/repos"
	"github.com/do87/poly/src/api/handlers/agents/usecases"
	"github.com/do87/poly/src/db"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	route  *chi.Mux
	repo   *repos.Repo
	uc     *usecases.Usecase
	agents *agents
	keys   *keys
}

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
			uc: usecases.NewAgentKeysUsecase(repo.Keys),
		},
	}
	p.setRoutes()
}
