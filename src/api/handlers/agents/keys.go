package agents

import (
	"net/http"

	"github.com/do87/poly/src/api/handlers/agents/present"
	"github.com/do87/poly/src/api/handlers/agents/usecases"
	"github.com/go-chi/render"
)

type keys struct {
	uc *usecases.Usecase
}

func (k *keys) list(u *usecases.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := u.Keys.List(r.Context(), r)
		if err != nil {
			render.JSON(w, r, present.Error(w, r, http.StatusInternalServerError, err))
			return
		}
		render.JSON(w, r, present.Keys(data))
	}
}
