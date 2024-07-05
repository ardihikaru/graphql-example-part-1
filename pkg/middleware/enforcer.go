package middleware

import (
	"net/http"
)

// AuthorizeResolver is a middleware to manage resolver function's access control based on the captured session
//
//	this method should be called after SessionCtx()
func (rs *Resource) AuthorizeResolver(act string) func(next http.Handler) http.Handler {
	return rs.utility.AuthorizeResolver(act)
}
