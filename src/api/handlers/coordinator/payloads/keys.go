package payloads

import (
	"time"

	"github.com/do87/poly/src/api/handlers/coordinator/models"
)

// KeyCreate is the payload for creating agent key
type KeyCreate struct {
	UUID      string    `json:"id"`
	Name      string    `json:"name"`
	PublicKey string    `json:"public_key"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ToModel converts payload to model
func (k *KeyCreate) ToModel() models.Key {
	return models.Key{
		UUID:      k.UUID,
		Name:      k.Name,
		PublicKey: []byte(k.PublicKey),
		ExpiresAt: k.ExpiresAt,
	}
}
