package gqlrouter

import (
	"github.com/go-chi/chi"

	"github.com/ardihikaru/graphql-example-part-1/internal/application"
	"github.com/ardihikaru/graphql-example-part-1/internal/router/graphql/role"
	"github.com/ardihikaru/graphql-example-part-1/internal/router/graphql/user"
	"github.com/ardihikaru/graphql-example-part-1/internal/service/session"
	"github.com/ardihikaru/graphql-example-part-1/internal/storage/resourcerolemap"

	"github.com/ardihikaru/graphql-example-part-1/pkg/jwtauth"
	"github.com/ardihikaru/graphql-example-part-1/pkg/middleware"
	"github.com/ardihikaru/graphql-example-part-1/pkg/middleware/middlewareutility"
)

func GraphQLRoutes(deps *application.Dependencies) chi.Router {
	router := chi.NewRouter()

	// builds resource group storage
	rsRoleStorage := &resourcerolemap.Storage{Db: deps.Db, Log: deps.Log}

	// initializes session middleware resource
	mwUtilSvc := middlewareutility.NewService(deps.Log, deps.Enforcer, rsRoleStorage)
	sessionSvc := session.NewService(deps.Log)
	mw := middleware.NewMiddleware(mwUtilSvc, sessionSvc)

	// adds middlewares to the graphql
	router.Use(mw.GraphQueryReader(deps.Cfg.GraphQL.PublicFunctions))
	router.Use(mw.Timeout(deps.Cfg.Http.WriteTimeout)) // returns 504
	router.Use(jwtauth.Verifier(deps.TokenAuth))       // Seeks, verifies and validates JWT tokens
	router.Use(mw.Authenticator())
	router.Use(mw.SessionCtx)               // extracts the session in the context
	router.Use(mw.AuthorizeResolver("use")) // adds policy to enable/disable to use the designated resolver function

	buildGraphQLTree(router, deps)

	return router
}

// buildGraphQLTree handles graphQL related route(s)
func buildGraphQLTree(r *chi.Mux, deps *application.Dependencies) {
	r.Mount(usergraph.BaseRoute, usergraph.BuildRouter(deps))
	r.Mount(rolegraph.BaseRoute, rolegraph.BuildRouter(deps))
}
