package graph

import (
	"github.com/developersismedika/sqlx"

	"github.com/ardihikaru/graphql-example-part-1/pkg/config"
	"github.com/ardihikaru/graphql-example-part-1/pkg/enforcer"
	"github.com/ardihikaru/graphql-example-part-1/pkg/jwtauth"
	"github.com/ardihikaru/graphql-example-part-1/pkg/logger"
	"github.com/ardihikaru/graphql-example-part-1/pkg/service/auth"
	"github.com/ardihikaru/graphql-example-part-1/pkg/service/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// utility provides the interface for the functionality of user resolver
type utility interface {
	Authorize(userId int, username string) (*auth.Token, error)
}

type Resolver struct {
	// inherits all of its functions
	*user.Service

	utility utility

	// supporting variables
	cfg       *config.Config
	db        *sqlx.DB
	log       *logger.Logger
	tokenAuth *jwtauth.JWTAuth
	enforcer  *enforcer.Enforcer
}

// NewResolver creates a new resolver
func NewResolver(cfg *config.Config, db *sqlx.DB, log *logger.Logger, tokenAuth *jwtauth.JWTAuth,
	enforcer *enforcer.Enforcer, userSvc *user.Service, utility utility) *Resolver {
	return &Resolver{
		Service:   userSvc,
		cfg:       cfg,
		db:        db,
		log:       log,
		tokenAuth: tokenAuth,
		enforcer:  enforcer,
		utility:   utility,
	}
}
