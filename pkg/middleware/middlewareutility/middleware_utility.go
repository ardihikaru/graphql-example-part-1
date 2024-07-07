package middlewareutility

import (
	"fmt"
	"net/http"

	"go.uber.org/zap/zapcore"

	"github.com/ardihikaru/graphql-example-part-1/pkg/enforcer"
	"github.com/ardihikaru/graphql-example-part-1/pkg/logger"
	"github.com/ardihikaru/graphql-example-part-1/pkg/middleware"
	"github.com/ardihikaru/graphql-example-part-1/pkg/service/session"
)

// storage provides the interface for the functionality to fetch information from the table resource_group in the DB
type storage interface {
	GetObjListOwner(resource, role string) ([]string, error)
}

type Service struct {
	*enforcer.Enforcer
	log     *logger.Logger
	storage storage
}

// NewService creates a new timeout handler service
func NewService(log *logger.Logger, enforcerPolicy *enforcer.Enforcer, storage storage) *Service {
	return &Service{
		log:      log,
		Enforcer: enforcerPolicy,
		storage:  storage,
	}
}

// Log logs message based on the log level
func (svc *Service) Log(level zapcore.Level, msg string) {
	switch level {
	case zapcore.DebugLevel:
		svc.log.Debug(msg)
	case zapcore.InfoLevel:
		svc.log.Info(msg)
	case zapcore.WarnLevel:
		svc.log.Warn(msg)
	case zapcore.ErrorLevel:
		svc.log.Error(msg)
	case zapcore.FatalLevel:
		svc.log.Fatal(msg)
	case zapcore.PanicLevel:
		svc.log.Panic(msg)
	default:
		svc.log.Error(fmt.Sprintf("unexpected log level type: %d", level))
	}
}

// AuthorizeResolver is a middleware to manage resolver function's access control based on the captured session
func (svc *Service) AuthorizeResolver(act string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var err error
			var ok bool

			// any public function will not be restricted to access
			isPublicFunction := r.Context().Value(middleware.PublicFunctionKey).(bool)
			if isPublicFunction {
				next.ServeHTTP(w, r)
				return
			}

			// gets current user/subject
			sessionData := r.Context().Value(middleware.SessionKey).(session.Session)
			svc.log.Debug(fmt.Sprintf("captured subject (userId/username): %d/%s", sessionData.UserId,
				sessionData.Username))
			sub := sessionData.UserId

			// loads policy from Database
			err = svc.LoadPolicy()
			if err != nil {
				svc.log.Warn("failed to load enforcer policy")
				http.Error(w, http.StatusText(http.StatusPreconditionFailed), http.StatusPreconditionFailed)

				return
			}

			// Casbin enforces policy for each object, whether the subject has an access or not
			functionName := r.Context().Value(middleware.RequestFunctionNameKey).(string)
			ok, err = svc.Enforce(fmt.Sprint(sub), functionName, act)
			if err != nil {
				svc.log.Warn("error occurred when authorizing user")
				http.Error(w, http.StatusText(http.StatusPreconditionFailed), http.StatusPreconditionFailed)

				return
			}
			if !ok {
				svc.log.Warn(fmt.Sprintf("this user (%s / %v) is not authorized",
					sessionData.Username, sessionData.UserId))
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
