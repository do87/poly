package payloads

import (
	"time"

	"github.com/do87/poly/src/api/handlers/mesh/models"
	"github.com/docker/distribution/uuid"
)

// KeyCreate is the payload for creating agent key
type KeyCreate struct {
	Name      string    `json:"name"`
	PublicKey string    `json:"public_key"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// ToModel converts payload to model
func (k KeyCreate) ToModel() models.Key {
	if k.ExpiresAt.IsZero() {
		k.ExpiresAt = time.Now().Add(time.Hour * 24 * 365)
	}
	return models.Key{
		UUID:      uuid.Generate().String(),
		Name:      k.Name,
		PublicKey: []byte(k.PublicKey),
		ExpiresAt: k.ExpiresAt,
	}
}
