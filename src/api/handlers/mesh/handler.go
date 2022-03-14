package mesh

import (
	"context"

	"github.com/do87/poly/src/api/handlers/mesh/repos"
	"github.com/do87/poly/src/api/handlers/mesh/usecases"
	"github.com/do87/poly/src/auth"
	"github.com/do87/poly/src/db"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	route  *chi.Mux
	repo   *repos.Repo
	auth   *auth.General
	agents *agents
	keys   *keys
}

// Handler handles agent related routes and functionality
func Handler(r *chi.Mux, d *db.DB) {
	repo := repos.New(d)
	a := &auth.General{}
	p := &handler{
		route: r,
		repo:  repo,
		auth:  a,

		// Link usecases:
		agents: &agents{
			uc:   usecases.NewAgentsUsecase(repo.Agents),
			auth: a,
		},
		keys: &keys{
			uc: usecases.NewKeysUsecase(repo.Keys),
		},
	}
	p.handleAuth().setRoutes()
}

func (h *handler) handleAuth() *handler {
	if h.auth.KeyExists() {
		return h
	}
	key, err := h.repo.Keys.FirstOrCreateGeneralKey(context.Background(), h.auth.GenerateKey())
	if err != nil {
		panic(err)
	}
	h.auth.SetKey(key.PublicKey)
	return h
}
