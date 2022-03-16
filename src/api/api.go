package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/do87/poly/src/db"
	"github.com/do87/poly/src/logger"
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
	DBConn   string
}

// cleanup is a type of a function to defer
type cleanup func() error
type API struct {
	router *chi.Mux
	db     *db.DB
	log    *logger.Logger
	config Config
}

// New creates a new API
func New(log *logger.Logger, db *db.DB, c Config) *API {
	return &API{
		log:    log,
		router: newRouter(log),
		config: c,
		db:     db,
	}
}

// Register registers the API
func (a *API) Register(handlers ...Handler) *API {
	for _, handler := range handlers {
		handler(a.router, a.db)
	}
	return a
}

// Run runs the API
func (a *API) Run() {
	a.log.Info("running server...", "host", a.serverStr())
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

func newRouter(log *logger.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Timeout(5 * time.Second))
	r.Use(middleware.RedirectSlashes)
	r.Use(log.ChiMiddleware())
	r.Use(render.SetContentType(render.ContentTypeJSON))
	return r
}
