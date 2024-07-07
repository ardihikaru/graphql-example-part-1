package graph

import (
	"crypto/rsa"
	"github.com/developersismedika/sqlx"

	"github.com/ardihikaru/graphql-example-part-1/pkg/config"
	"github.com/ardihikaru/graphql-example-part-1/pkg/enforcer"
	"github.com/ardihikaru/graphql-example-part-1/pkg/graph/user/model"
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

// EncryptPassword encrypts password
func (r *Resolver) EncryptPassword(password string, publicKeyRSA *rsa.PublicKey) ([]byte, error) {
	return r.EncryptPassword(password, publicKeyRSA)
}

// Create creates a new user
func (r *Resolver) Create(data model.UserInput, CreateUserID *int64) (*model.User, error) {
	return r.Create(data, CreateUserID)
}

// GetUserIdByUsername gets user data by username
func (r *Resolver) GetUserIdByUsername(username string) (int, error) {
	return r.GetUserIdByUsername(username)
}

// Authenticate authenticates the provided credential
func (r *Resolver) Authenticate(userName string, password string) (bool, error) {
	return r.Authenticate(userName, password)
}

// List gets list of user
func (r *Resolver) List(userIdStr, statusCd string) ([]*model.User, error) {
	return r.List(userIdStr, statusCd)
}

// GetById gets user data by ID
func (r *Resolver) GetById(userId int64) (*model.User, error) {
	return r.GetById(userId)
}
