package models

import (
	"time"

	"github.com/lib/pq"
)

type Agent struct {
	UUID      string         `gorm:"primaryKey"`
	Hostname  string         `gorm:"primaryKey"`
	Labels    pq.StringArray `gorm:"type:text[]"`
	Plans     pq.StringArray `gorm:"type:text[]"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName to ensure correct table naming
func (Agent) TableName() string {
	return "agents"
}
