package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"

	"github.com/ardihikaru/graphql-example-part-1/pkg/graph/role/generated"
	"github.com/ardihikaru/graphql-example-part-1/pkg/graph/role/model"
)

// RoleCreate is the resolver for the roleCreate field.
func (r *mutationResolver) RoleCreate(ctx context.Context, data model.RoleInput) (*model.Role, error) {
	panic(fmt.Errorf("not implemented: RoleCreate - roleCreate"))
}

// RoleGet is the resolver for the roleGet field.
func (r *queryResolver) RoleGet(ctx context.Context, roleID *string) (*model.Role, error) {
	return &model.Role{
		RoleID: "123",
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
