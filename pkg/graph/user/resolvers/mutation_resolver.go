package graph

import (
	"crypto/rsa"

	"github.com/ardihikaru/graphql-example-part-1/pkg/graph/user/model"
)

// EncryptPassword encrypts password
func (r *mutationResolver) EncryptPassword(password string, publicKeyRSA *rsa.PublicKey) ([]byte, error) {
	return r.Resolver.EncryptPassword(password, publicKeyRSA)
}

// Create creates a new user
func (r *mutationResolver) Create(data model.UserInput, CreateUserID *int64) (*model.User, error) {
	return r.Resolver.Create(data, CreateUserID)
}

// GetUserIdByUsername gets user data by username
func (r *mutationResolver) GetUserIdByUsername(username string) (int, error) {
	return r.Resolver.GetUserIdByUsername(username)
}

// Authenticate authenticates the provided credential
func (r *mutationResolver) Authenticate(userName string, password string) (bool, error) {
	return r.Resolver.Authenticate(userName, password)
}
