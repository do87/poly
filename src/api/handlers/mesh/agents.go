package mesh

import (
	"net/http"

	"github.com/do87/poly/src/api/handlers/mesh/present"
	"github.com/do87/poly/src/api/handlers/mesh/repos"
	"github.com/do87/poly/src/api/handlers/mesh/usecases"
	"github.com/do87/poly/src/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type agents struct {
	uc   *usecases.Usecase
	auth *auth.General
}

func newAgentHandler(repo *repos.Repo, general *auth.General) *agents {
	return &agents{
		uc:   usecases.NewAgentsUsecase(repo.Agents),
		auth: general,
	}
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

func (a *agents) register(agentsUc, keysUc *usecases.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		agent, err := agentsUc.Agents.Register(r.Context(), r, keysUc.Keys)
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
