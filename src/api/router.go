package api

import (
	"time"

	"github.com/do87/poly/src/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func NewRouter(log *logger.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Timeout(5 * time.Second))
	r.Use(middleware.RedirectSlashes)
	r.Use(log.ChiMiddleware())
	r.Use(render.SetContentType(render.ContentTypeJSON))
	return r
}
