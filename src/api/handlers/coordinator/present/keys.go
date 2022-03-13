package present

import (
	"fmt"
	"time"

	"github.com/amalfra/etag"
	"github.com/do87/poly/src/api/handlers/coordinator/models"
)

// key consts
const (
	KEY_AGENT_KEY  = "poly:agent-key"
	KEY_AGENT_KEYS = "poly:agent-keys"
)

type key struct {
	UUID      string    `json:"id"`
	Name      string    `json:"name"`
	PublicKey string    `json:"public_key"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FromModel converts model to presenter struct
func (key) FromModel(m models.Key) key {
	return key{
		UUID:      m.UUID,
		Name:      m.Name,
		PublicKey: string(m.PublicKey),
		ExpiresAt: m.ExpiresAt,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// Keys presents agent keys
func Keys(keyModels []models.Key) Presentor {
	p := make([]key, 0)
	a := key{}
	u := ""
	for _, m := range keyModels {
		u += fmt.Sprintf("%s-%s-%s;", m.UUID, m.Name, m.UpdatedAt.String())
		p = append(p, a.FromModel(m))
	}
	return wrap(KEY_AGENT_KEYS, etag.Generate(u, true), p)
}

// Key presents an agent key
func Key(m models.Key) Presentor {
	a := key{}
	u := fmt.Sprintf("%s-%s-%s;", m.UUID, m.Name, m.UpdatedAt.String())
	return wrap(KEY_AGENT_KEY, etag.Generate(u, true), a.FromModel(m))
}
