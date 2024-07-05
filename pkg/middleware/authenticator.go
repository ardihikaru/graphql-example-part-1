package middleware

import (
	"context"
	"net/http"

	"github.com/ardihikaru/graphql-example-part-1/pkg/jwtauth"
)

func validateToken(ctx context.Context, w http.ResponseWriter) bool {
	token, _, err := jwtauth.FromContext(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return false
	}

	if token == nil || jwtauth.ValidateToken(token) != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return false
	}

	return true
}

func Authenticator() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isPublicFunction := r.Context().Value(PublicFunctionKey).(bool)

			// Allow unauthenticated users in
			if isPublicFunction {
				next.ServeHTTP(w, r)
				return
			}

			validToken := validateToken(r.Context(), w)
			if !validToken {
				// FYI: http header has been added inside validateToken() function

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
