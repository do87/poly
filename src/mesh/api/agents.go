package mesh

import (
	"net/http"

	"github.com/do87/poly/src/mesh/api/present"
	"github.com/do87/poly/src/mesh/api/usecases"
	"github.com/do87/poly/src/mesh/repos"
	"github.com/do87/poly/src/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type agents struct {
	subHandler
	auth *auth.General
}

func newAgentHandler(repo *repos.Repo, general *auth.General) *agents {
	a := &agents{
		auth: general,
	}
	a.uc = usecases.NewAgentsUsecase(repo.Agents)
	return a
}

func (a *agents) list(u *usecases.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := u.Agents.List(r.Context(), r)
		if err != nil {
			render.JSON(w, r, present.Error(w, r, http.StatusInternalServerError, err))
			return
		}
		render.JSON(w, r, present.Agents(data))
	}
}

func (a *agents) register(u *usecases.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		agent, err := u.Agents.Register(r.Context(), r, u.Keys)
		if err != nil {
			render.JSON(w, r, present.Error(w, r, http.StatusInternalServerError, err))
			return
		}
		token, err := a.auth.Token(agent.Hostname)
		if err != nil {
			render.JSON(w, r, present.Error(w, r, http.StatusInternalServerError, err))
			return
		}
		_ = token
		render.JSON(w, r, present.AccessToken(token))
	}
}

func (a *agents) deregister(u *usecases.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		a, err := u.Agents.Deregister(r.Context(), r, id)
		if err != nil {
			render.JSON(w, r, present.Error(w, r, http.StatusInternalServerError, err))
			return
		}
		render.JSON(w, r, present.Generic(present.KEY_AGENT, "", a))
	}
}
