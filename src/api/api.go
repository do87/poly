package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/do87/poly/src/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Handler func(*chi.Mux)

type Config struct {
	BindAddr string
	BindPort int
}

type api struct {
	router *chi.Mux
	config Config
}

func New(c Config) *api {
	return &api{
		router: newChiRouter(),
		config: c,
	}
}

func (a *api) Register(handlers ...Handler) *api {
	for _, handler := range handlers {
		handler(a.router)
	}
	return a
}

func (a *api) Run() {
	log, logsync := logger.New()
	defer logsync()
	a.router.Use(log.ChiMiddleware())

	if err := http.ListenAndServe(a.serverStr(), a.router); err != nil {
		panic(err)
	}
}

func (a *api) serverStr() string {
	addr := a.config.BindAddr
	if addr == "" {
		addr = "0.0.0.0"
	}

	port := a.config.BindPort
	if port == 0 {
		port = 8080
	}

	return fmt.Sprintf("%s:%d", addr, port)
}

func newChiRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Timeout(5 * time.Second))
	r.Use(middleware.RedirectSlashes)
	// r.Use(log.ChiMiddleware())
	r.Use(render.SetContentType(render.ContentTypeJSON))
	return r
}
