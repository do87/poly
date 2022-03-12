package agents

import (
	"net/http"

	"github.com/do87/poly/src/api/handlers/agents/present"
	"github.com/do87/poly/src/api/handlers/agents/usecases"
	"github.com/go-chi/render"
)

type agents struct {
	uc *usecases.Usecase
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
		render.JSON(w, r, nil)
	}
}