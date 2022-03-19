package middlewares

import (
	"net/http"

	"github.com/do87/poly/src/mesh/common"
	"github.com/do87/poly/src/pkg/auth"
)

// VerifyAgent is a middleware to verify an agent's access token
func VerifyAgent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuid, err := auth.ValidateGeneralTokenHeader(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
		r = common.AddToContext(r, common.ContextAgentUUID, uuid)
		next.ServeHTTP(w, r)
	})
}
