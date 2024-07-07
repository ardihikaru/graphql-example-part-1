package rolegraph

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"

	"github.com/ardihikaru/graphql-example-part-1/internal/application"

	"github.com/ardihikaru/graphql-example-part-1/pkg/graph/role/generated"
	graph "github.com/ardihikaru/graphql-example-part-1/pkg/graph/role/resolvers"
)

const (
	BaseRoute = "/role"
)

// BuildRouter builds the router
func BuildRouter(deps *application.Dependencies) http.Handler {
	r := chi.NewRouter()

	buildGraphQLHandler(r, deps)

	return r
}

// buildGraphQLHandler builds graphQL handler
func buildGraphQLHandler(r *chi.Mux, deps *application.Dependencies) {
	// creates GraphQL server
	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: graph.NewResolver(deps.Cfg, deps.Db, deps.Log, deps.TokenAuth, deps.Enforcer),
	}))

	srv.AddTransport(transport.POST{})
	srv.Use(extension.Introspection{})

	// registers GraphQL playground endpoint
	r.Handle("/", playground.Handler("GraphQL playground", fmt.Sprintf("%s/query", BaseRoute)))

	// registers GraphQL endpoint
	r.Handle("/query", srv)
}
