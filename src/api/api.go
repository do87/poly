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

type cleanup func() error
type api struct {
	router  *chi.Mux
	log     *logger.Logger
	config  Config
	cleanup []cleanup
}

func New(c Config) *api {
	log, logsync := logger.New()

	return &api{
		log:     log,
		router:  newRouter(log),
		config:  c,
		cleanup: []cleanup{logsync},
	}
}

func (a *api) Register(handlers ...Handler) *api {
	for _, handler := range handlers {
		handler(a.router)
	}
	return a
}

func (a *api) Run() {
	defer a.Cleanup()
	a.log.Info("running server...", "host", a.serverStr())
	if err := http.ListenAndServe(a.serverStr(), a.router); err != nil {
		panic(err)
	}
}

func (a *api) Cleanup() {
	for _, c := range a.cleanup {
		c()
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

func newRouter(log *logger.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Timeout(5 * time.Second))
	r.Use(middleware.RedirectSlashes)
	r.Use(log.ChiMiddleware())
	r.Use(render.SetContentType(render.ContentTypeJSON))
	return r
}
