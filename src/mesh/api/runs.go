package mesh

import (
	"net/http"

	"github.com/do87/poly/src/mesh/api/present"
	"github.com/do87/poly/src/mesh/api/usecases"
	"github.com/do87/poly/src/mesh/repos"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type runs subHandler

func newRunsHandler(repo *repos.Repo) *runs {
	return &runs{
		uc: usecases.NewRunsUsecase(repo.Runs),
	}
}

func (k *runs) list(u *usecases.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := u.Runs.List(r.Context(), r)
		if err != nil {
			render.JSON(w, r, present.Error(w, r, http.StatusInternalServerError, err))
			return
		}
		render.JSON(w, r, present.Runs(data))
	}
}

func (k *runs) create(u *usecases.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := u.Runs.Create(r.Context(), r)
		if err != nil {
			render.JSON(w, r, present.Error(w, r, http.StatusInternalServerError, err))
			return
		}
		render.JSON(w, r, present.Run(data))
	}
}

func (k *runs) update(u *usecases.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		data, err := u.Runs.Update(r.Context(), r, id)
		if err != nil {
			render.JSON(w, r, present.Error(w, r, http.StatusInternalServerError, err))
			return
		}
		render.JSON(w, r, present.Run(data))
	}
}
