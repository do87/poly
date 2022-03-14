package ctx

import (
	"context"
	"net/http"
)

// MeshCtxKey is context keys type for mesh API
type MeshCtxKey string

// CtxHostname is the value the key that will be added to the context
const CtxAgentHostname MeshCtxKey = "hostname"

func Add(r *http.Request, key MeshCtxKey, value string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key, value))
}
