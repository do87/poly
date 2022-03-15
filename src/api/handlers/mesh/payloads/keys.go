package payloads

import (
	"time"

	"github.com/do87/poly/src/api/handlers/mesh/models"
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
		Name:      k.Name,
		PublicKey: []byte(k.PublicKey),
		ExpiresAt: k.ExpiresAt,
	}
}

// EncodedKey is the key related payload used during agent registration
type EncodedKey struct {
	Name    string `json:"name"`
	Encoded string `json:"key"`
}
