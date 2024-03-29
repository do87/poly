package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/do87/poly/src/pkg/db"
	"github.com/do87/poly/src/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// Handler defines routes for a service
type Handler func(*chi.Mux, *db.DB)

// Config is the API config
type Config struct {
	BindAddr string
	BindPort int
}

// cleanup is a type of a function to defer
type cleanup func() error

// API is the API service
type API struct {
	router *chi.Mux
	db     *db.DB
	log    logger.Log
	config Config
}

// New creates a new API
func New(log logger.Log, db *db.DB, c Config) *API {
	return &API{
		log:    log,
		router: newRouter(log),
		config: c,
		db:     db,
	}
}

// Register registers API routes using handler functions
func (a *API) Register(handlers ...Handler) *API {
	for _, handler := range handlers {
		handler(a.router, a.db)
	}
	return a
}

// Run runs the API
func (a *API) Run() {
	a.log.Info(fmt.Sprintf("🚀 running API server on %s", a.serverStr()))
	if err := http.ListenAndServe(a.serverStr(), a.router); err != nil {
		panic(err)
	}
}

func (a *API) serverStr() string {
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

func newRouter(log logger.Log) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Timeout(5 * time.Second))
	r.Use(middleware.RedirectSlashes)
	r.Use(log.ChiMiddleware())
	r.Use(render.SetContentType(render.ContentTypeJSON))
	return r
}
