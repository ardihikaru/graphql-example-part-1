package session

import (
	"net/http"

	"github.com/ardihikaru/graphql-example-part-1/pkg/logger"
	"github.com/ardihikaru/graphql-example-part-1/pkg/service/session"
)

type Service struct {
	*session.Service
}

// NewService creates a new timeout handler service
func NewService(log *logger.Logger) *Service {
	service := session.NewService(log)

	return &Service{
		Service: service,
	}
}

// SessionCtx enriches the request with the captured JWT private claims
//
//	this function inherits the implementation from its parent
func (svc *Service) SessionCtx(next http.Handler) http.Handler {
	return svc.Service.SessionCtx(next)
}
