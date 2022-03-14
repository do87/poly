package ctx

import (
	"context"
	"net/http"
)

// MeshCtxKey is context keys type for mesh API
type MeshCtxKey string

// CtxAgentHostname is the key that will be added to the context
const CtxAgentHostname MeshCtxKey = "hostname"

// Add adds a key of type MeshCtxKey and a provided string value to the request context
func Add(r *http.Request, key MeshCtxKey, value string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key, value))
}
