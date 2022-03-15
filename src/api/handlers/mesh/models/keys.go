package models

import (
	"time"
)

// Key model
type Key struct {
	Name      string `gorm:"primaryKey"`
	PublicKey []byte
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

// TableName to ensure correct table naming
func (Key) TableName() string {
	return "keys"
}
