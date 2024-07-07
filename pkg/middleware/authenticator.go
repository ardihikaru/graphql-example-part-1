package middleware

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap/zapcore"

	"github.com/ardihikaru/graphql-example-part-1/pkg/jwtauth"
)

func (rs *Resource) validateToken(ctx context.Context, w http.ResponseWriter) bool {
	token, _, err := jwtauth.FromContext(ctx)

	if err != nil {
		rs.utility.Log(zapcore.ErrorLevel, fmt.Sprintf("failed to fetch token from header: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return false
	}

	if token == nil || jwtauth.ValidateToken(token) != nil {
		rs.utility.Log(zapcore.ErrorLevel, "invalid token in the header")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return false
	}

	return true
}

func (rs *Resource) Authenticator() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// do nothing on any GET method request
			if r.Method == http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			isPublicFunction := r.Context().Value(PublicFunctionKey).(bool)

			// Allow unauthenticated users in
			if isPublicFunction {
				next.ServeHTTP(w, r)
				return
			}

			validToken := rs.validateToken(r.Context(), w)
			if !validToken {
				// FYI: http header has been added inside validateToken() function

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
