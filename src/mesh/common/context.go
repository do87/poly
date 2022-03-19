package common

import (
	"context"
	"net/http"
)

// MeshContextVar is context keys type for mesh API
type MeshContextVar string

// ContextAgentUUID is the key that will be added to the context
const ContextAgentUUID MeshContextVar = "uuid"

// AddToContext adds a key of type MeshCtxKey and a provided string value to the request context
func AddToContext(r *http.Request, key MeshContextVar, value string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key, value))
}
