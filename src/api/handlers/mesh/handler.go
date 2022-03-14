package mesh

import (
	"context"
	"os"

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
	p.setGlobalKey().setRoutes()
}

func (h *handler) setGlobalKey() *handler {
	meshKey := os.Getenv(usecases.MeshGlobalKey)
	if meshKey != "" {
		return h
	}
	key, err := h.keys.uc.Keys.CreateGlobalKeyIfNotExists(context.Background())
	if err != nil {
		panic(err)
	}
	os.Setenv(usecases.MeshGlobalKey, string(key.PublicKey))
	return h
}
