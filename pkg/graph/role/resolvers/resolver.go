package graph

import (
	"github.com/developersismedika/sqlx"

	"github.com/ardihikaru/graphql-example-part-1/pkg/config"
	"github.com/ardihikaru/graphql-example-part-1/pkg/enforcer"
	"github.com/ardihikaru/graphql-example-part-1/pkg/jwtauth"
	"github.com/ardihikaru/graphql-example-part-1/pkg/logger"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	// supporting variables
	cfg       *config.Config
	db        *sqlx.DB
	log       *logger.Logger
	tokenAuth *jwtauth.JWTAuth
	enforcer  *enforcer.Enforcer
}

// NewResolver creates a new resolver
func NewResolver(cfg *config.Config, db *sqlx.DB, log *logger.Logger, tokenAuth *jwtauth.JWTAuth,
	enforcer *enforcer.Enforcer) *Resolver {
	return &Resolver{
		cfg:       cfg,
		log:       log,
		db:        db,
		tokenAuth: tokenAuth,
		enforcer:  enforcer,
	}
}
