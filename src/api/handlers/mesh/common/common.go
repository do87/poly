package common

import (
	"context"
	"net/http"

	"github.com/do87/poly/src/api/handlers/mesh/models"
	"github.com/do87/poly/src/api/handlers/mesh/payloads"
	"github.com/do87/poly/src/auth"
)

// MeshContextType is context keys type for mesh API
type MeshContextVar string

// ContextAgentHostname is the key that will be added to the context
const ContextAgentHostname MeshContextVar = "hostname"

// AddToContext adds a key of type MeshCtxKey and a provided string value to the request context
func AddToContext(r *http.Request, key MeshContextVar, value string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key, value))
}

// ProcessRegisterKey processes a registration key using the provided public key
func ProcessRegisterKey(key models.Key, payload payloads.AgentRegister) error {
	meshKey := auth.MeshRegisterKey{
		Name:      key.Name,
		PublicKey: key.PublicKey,
	}
	t, err := meshKey.Decode(payload.EncodedKey.Encoded)
	if err != nil {
		return err
	}
	if err := auth.ValidateRegisterToken(t, payload.Hostname); err != nil {
		return err
	}
	return nil
}
