package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is a wrapper around gorm
type DB struct {
	db *gorm.DB
}

// New returns a new database service
func New(dialector gorm.Dialector) (*DB, error) {
	db, err := gorm.Open(dialector)
	if err != nil {
		return nil, fmt.Errorf("error opening to DB: %s", err)
	}

	return &DB{
		db: db,
	}, nil
}

// NewPostgres returns a new postgres database service from connection string
func NewPostgres(conn string) (*DB, error) {
	return New(postgres.Open(conn))
}

// GetDB returns a pointer to the DB
func (d *DB) GetDB() *gorm.DB {
	return d.db
}

// Close closes connection to the DB
func (d *DB) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Migrate runs table migrations
func (d *DB) Migrate(tables ...interface{}) error {
	return d.db.AutoMigrate(tables...)
}
