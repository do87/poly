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
	route *chi.Mux
	repo  *repos.Repo
	auth  *auth.General

	// sub handlers:
	agents *agents
	keys   *keys
	runs   *runs
}

type subHandler struct {
	uc *usecases.Usecase
}

// Handler handles agent related routes and functionality
func Handler(r *chi.Mux, d *db.DB) {
	repo := repos.New(d)
	general := &auth.General{}
	p := &handler{
		route: r,
		repo:  repo,
		auth:  general,

		// sub handlers:
		agents: newAgentHandler(repo, general),
		keys:   newKeysHandler(repo),
		runs:   newRunsHandler(repo),
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
