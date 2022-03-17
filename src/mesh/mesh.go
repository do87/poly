package mesh

import (
	meshAPI "github.com/do87/poly/src/mesh/api"
	"github.com/do87/poly/src/mesh/routines"
	"github.com/do87/poly/src/pkg/api"
	"github.com/do87/poly/src/pkg/db"
	"github.com/do87/poly/src/pkg/health"
	"github.com/do87/poly/src/pkg/logger"
)

type mesh struct {
	config  Config
	log     *logger.Logger
	db      *db.DB
	cleanup []cleanup
	api     *api.API
}

// APIConfig mirrors the API config
type APIConfig = api.Config

// Config is the mesh server config
type Config struct {
	API    APIConfig
	DBConn string
}

// cleanup is a type of a function to defer
type cleanup func() error

// New creates a new mesh server
func New(c Config) *mesh {
	log, logsync := logger.New()
	m := &mesh{
		log:     log,
		cleanup: []cleanup{logsync},
		db:      setupDatabase(c),
	}
	m.api = api.New(log, m.db, c.API)
	m.cleanup = append(m.cleanup, m.db.Close)
	m.Register(
		health.Handler,
		meshAPI.Handler,
	)
	return m
}

func setupDatabase(c Config) *db.DB {
	if c.DBConn == "" {
		return nil
	}
	db, err := db.NewPostgres(c.DBConn)
	if err != nil {
		panic(err)
	}
	return db
}

// Register registers APIs
func (m *mesh) Register(handlers ...api.Handler) *mesh {
	m.api.Register(handlers...)
	return m
}

// Run runs the mesh server
func (m *mesh) Run() {
	defer m.Cleanup()
	go routines.Run(m.db, m.log)
	m.api.Run()
}

// Cleanup runs cleanup functions
func (m *mesh) Cleanup() {
	for _, c := range m.cleanup {
		c()
	}
}
