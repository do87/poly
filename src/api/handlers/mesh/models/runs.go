package models

import (
	"time"

	"github.com/lib/pq"
)

// Run model
type Run struct {
	UUID      string `gorm:"primaryKey"`
	AgentUUID string
	Plan      string `gorm:"not null;default:null"`
	Labels    pq.StringArray
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName to ensure correct table naming
func (Run) TableName() string {
	return "runs"
}
