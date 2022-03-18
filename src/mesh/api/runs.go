package mesh

import (
	"net/http"

	"github.com/do87/poly/src/mesh/api/present"
	"github.com/do87/poly/src/mesh/repos"
	"github.com/do87/poly/src/mesh/usecases"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type runs subHandler

func newRunsHandler(repo *repos.Repo) *runs {
	return &runs{
		uc: usecases.NewRunsUsecase(repo.Runs),
	}
}

func (h *runs) list(u *usecases.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := u.Runs.List(r.Context(), r)
		if err != nil {
			render.JSON(w, r, present.Error(w, r, http.StatusInternalServerError, err))
			return
		}
		render.JSON(w, r, present.Runs(data))
	}
}

func (h *runs) listPending(ruc, a *usecases.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if _, err := a.Agents.Ping(r.Context(), r, id); err != nil {
			render.JSON(w, r, present.Error(w, r, http.StatusInternalServerError, err))
			return
		}
		data, err := ruc.Runs.ListPending(r.Context(), r, id)
		if err != nil {
			render.JSON(w, r, present.Error(w, r, http.StatusInternalServerError, err))
			return
		}
		render.JSON(w, r, present.Runs(data))
	}
}

func (h *runs) create(u *usecases.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := u.Runs.Create(r.Context(), r)
		if err != nil {
			render.JSON(w, r, present.Error(w, r, http.StatusInternalServerError, err))
			return
		}
		render.JSON(w, r, present.Run(data))
	}
}

func (h *runs) update(u *usecases.Usecase) http.HandlerFunc {
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
