package models

import (
	"time"

	"github.com/lib/pq"
)

// Run model
type Run struct {
	UUID       string `gorm:"primaryKey"`
	Agent      string
	Plan       string `gorm:"not null;default:null"`
	Labels     pq.StringArray
	Status     string
	StartedAt  time.Time
	FinishedAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// TableName to ensure correct table naming
func (Run) TableName() string {
	return "runs"
}
