package auth

import (
	"github.com/ardihikaru/graphql-example-part-1/pkg/jwtauth"
	"github.com/ardihikaru/graphql-example-part-1/pkg/logger"
	"github.com/ardihikaru/graphql-example-part-1/pkg/service/auth"
)

// Service prepares the interfaces related with this auth service
type Service struct {
	*auth.Service
}

// NewService creates a new auth service
func NewService(log *logger.Logger, jwtExpTime int, tokenAuth *jwtauth.JWTAuth) *Service {
	service := auth.NewService(log, jwtExpTime, tokenAuth)

	return &Service{
		Service: service,
	}
}

// Authorize authorizes user credential
//
//	this function inherits the implementation from its parent
func (svc *Service) Authorize(userId int, username string) (*auth.Token, error) {
	return svc.Service.Authorize(userId, username)
}
