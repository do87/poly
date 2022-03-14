package middlewares

import (
	"net/http"

	"github.com/do87/poly/src/api/handlers/mesh/ctx"
	"github.com/do87/poly/src/auth"
)

func Agents(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hostname, err := auth.ValidateGeneralTokenHeader(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
		r = ctx.Add(r, ctx.CtxAgentHostname, hostname)
		next.ServeHTTP(w, r)
	})
}
