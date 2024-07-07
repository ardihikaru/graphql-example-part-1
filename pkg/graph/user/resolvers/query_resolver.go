package graph

import (
	"github.com/ardihikaru/graphql-example-part-1/pkg/graph/user/model"
)

// List gets list of user
func (r *queryResolver) List(userIdStr, statusCd string) ([]*model.User, error) {
	return r.Resolver.List(userIdStr, statusCd)
}

// GetById gets user data by ID
func (r *queryResolver) GetById(userId int64) (*model.User, error) {
	return r.Resolver.GetById(userId)
}
