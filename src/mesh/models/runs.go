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
	Payload    []byte `gorm:"default:null"`
	Status     string
	AssignedAt time.Time `gorm:"default:null"`
	StartedAt  time.Time `gorm:"default:null"`
	FinishedAt time.Time `gorm:"default:null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// TableName to ensure correct table naming
func (Run) TableName() string {
	return "runs"
}
