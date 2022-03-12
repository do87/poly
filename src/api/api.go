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
type api struct {
	router  *chi.Mux
	db      *db.DB
	log     *logger.Logger
	config  Config
	cleanup []cleanup
}

// New creates a new API
func New(c Config) *api {
	log, logsync := logger.New()
	api := &api{
		log:     log,
		router:  newRouter(log),
		config:  c,
		cleanup: []cleanup{logsync},
	}
	api.setupDatabase()
	return api
}

func (a *api) setupDatabase() {
	if a.config.DBConn == "" {
		return
	}
	db, err := db.NewPostgres(a.config.DBConn)
	if err != nil {
		panic(err)
	}
	a.db = db
	a.cleanup = append(a.cleanup, db.Close)
}

// Register registers the API
func (a *api) Register(handlers ...Handler) *api {
	for _, handler := range handlers {
		handler(a.router, a.db)
	}
	return a
}

// Run runs the API
func (a *api) Run() {
	defer a.Cleanup()
	a.log.Info("running server...", "host", a.serverStr())
	if err := http.ListenAndServe(a.serverStr(), a.router); err != nil {
		panic(err)
	}
}

// Cleanup run cleanup functions
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
