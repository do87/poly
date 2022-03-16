package mesh

import (
	"net/http"

	"github.com/do87/poly/src/api/handlers/mesh/present"
	"github.com/do87/poly/src/api/handlers/mesh/usecases"
	"github.com/do87/poly/src/mesh/repos"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type keys subHandler

func newKeysHandler(repo *repos.Repo) *keys {
	return &keys{
		uc: usecases.NewKeysUsecase(repo.Keys),
	}
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

func (k *keys) create(u *usecases.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := u.Keys.Create(r.Context(), r)
		if err != nil {
			render.JSON(w, r, present.Error(w, r, http.StatusInternalServerError, err))
			return
		}
		render.JSON(w, r, present.Key(data))
	}
}

func (k *keys) delete(u *usecases.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "name")
		if err := u.Keys.Delete(r.Context(), r, id); err != nil {
			render.JSON(w, r, present.Error(w, r, http.StatusInternalServerError, err))
			return
		}
		render.JSON(w, r, present.Generic(present.KEY_AGENT_KEY, "", nil))
	}
}
