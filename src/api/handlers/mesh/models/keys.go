package models

import (
	"time"
)

// Key model
type Key struct {
	UUID      string
	Name      string
	PublicKey []byte
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

// TableName to ensure correct table naming
func (Key) TableName() string {
	return "keys"
}
