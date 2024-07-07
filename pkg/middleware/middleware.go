package middleware

import (
	"net/http"

	"go.uber.org/zap/zapcore"
)

const (
	RequestId = "X-Request-Id"
)

// utility provides the interface for the functionality of logger.Logger and any other common utility
type utility interface {
	AuthorizeResolver(act string) func(next http.Handler) http.Handler
	Log(level zapcore.Level, msg string)
}

// session provides the interface for the functionality of session handler for the authentication
type session interface {
	SessionCtx(next http.Handler) http.Handler
}

type Resource struct {
	utility utility
	session session
}

// NewMiddleware creates a new middleware
func NewMiddleware(utility utility, session session) *Resource {
	return &Resource{
		utility: utility,
		session: session,
	}
}
