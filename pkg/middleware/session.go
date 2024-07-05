package middleware

import "net/http"

const (
	// SessionKey is the context key to store JWT private claims which is captured from the request
	SessionKey             = "session"
	PublicFunctionKey      = "public-function"
	RequestFunctionNameKey = "function-name"
)

// SessionCtx authorizes session credential
func (rs *Resource) SessionCtx(next http.Handler) http.Handler {
	return rs.session.SessionCtx(next)
}
