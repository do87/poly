package present

import (
	"fmt"
	"time"

	"github.com/amalfra/etag"
	"github.com/do87/poly/src/mesh/models"
)

// key constants
const (
	KEY_AGENT_KEY  = "poly:agent-key"
	KEY_AGENT_KEYS = "poly:agent-keys"
)

type TypeKey struct {
	Name      string    `json:"name"`
	PublicKey string    `json:"public_key"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FromModel converts model to presenter struct
func (TypeKey) FromModel(m models.Key) TypeKey {
	return TypeKey{
		Name:      m.Name,
		PublicKey: string(m.PublicKey),
		ExpiresAt: m.ExpiresAt,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

var _k = TypeKey{}

// Keys presents agent keys
func Keys(keyModels []models.Key) Presentor[[]TypeKey] {
	p := make([]TypeKey, 0)
	u := ""
	for _, m := range keyModels {
		u += fmt.Sprintf("%s-%s;", m.Name, m.UpdatedAt.String())
		p = append(p, _k.FromModel(m))
	}
	return wrap(KEY_AGENT_KEYS, etag.Generate(u, true), p)
}

// Key presents an agent key
func Key(m models.Key) Presentor[TypeKey] {
	u := fmt.Sprintf("%s-%s;", m.Name, m.UpdatedAt.String())
	return wrap(KEY_AGENT_KEY, etag.Generate(u, true), _k.FromModel(m))
}
