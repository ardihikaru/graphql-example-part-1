package middleware

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap/zapcore"
)

// Timeout is a middleware that cancels ctx after a given timeout and return
// a 504 Gateway Timeout error to the client.
//
// It's required that you select the ctx.Done() channel to check for the signal
// if the context has reached its deadline and return, otherwise the timeout
// signal will be just ignored.
//
// i.e. a route/handler may look like:
//
//	r.Get("/long", func(w http.ResponseWriter, r *http.Request) {
//		ctx := r.Context()
//		processTime := time.Duration(rand.Intn(4)+1) * time.Second
//
//		select {
//		case <-ctx.Done():
//			return
//
//		case <-time.After(processTime):
//			// The above channel simulates some hard work.
//		}
//
//		w.Write([]byte("done"))
//	})
//
// FYI: it adopted: https://github.com/go-chi/chi/blob/master/middleware/timeout.go
func (rs *Resource) Timeout(timeout time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer func() {
				cancel()
				if ctx.Err() == context.DeadlineExceeded {
					rs.utility.Log(zapcore.DebugLevel, "got a request timeout")
					w.WriteHeader(http.StatusGatewayTimeout)

					return
				}
			}()

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
