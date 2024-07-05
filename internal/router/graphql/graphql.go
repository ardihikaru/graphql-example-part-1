package gqlrouter

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"

	"github.com/ardihikaru/graphql-example-part-1/internal/application"
	"github.com/ardihikaru/graphql-example-part-1/internal/graph/generated"
	graph "github.com/ardihikaru/graphql-example-part-1/internal/graph/resolvers"
	"github.com/ardihikaru/graphql-example-part-1/internal/service/middlewareutility"
	"github.com/ardihikaru/graphql-example-part-1/internal/service/session"
	"github.com/ardihikaru/graphql-example-part-1/internal/storage/resourcerolemap"

	"github.com/ardihikaru/graphql-example-part-1/pkg/jwtauth"
	"github.com/ardihikaru/graphql-example-part-1/pkg/middleware"
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
	router.Use(middleware.GraphQueryReader(deps.Cfg.GraphQL.PublicFunctions))
	router.Use(mw.Timeout(deps.Cfg.Http.WriteTimeout)) // returns 504
	router.Use(jwtauth.Verifier(deps.TokenAuth))       // Seeks, verifies and validates JWT tokens
	router.Use(middleware.Authenticator())
	router.Use(mw.SessionCtx)               // extracts the session in the context
	router.Use(mw.AuthorizeResolver("use")) // adds policy to enable/disable to use the designated resolver function

	// Create GraphQL server
	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{
			SvcId:     deps.SvcId,
			Cfg:       deps.Cfg,
			Db:        deps.Db,
			Log:       deps.Log,
			TokenAuth: deps.TokenAuth,
			Enforcer:  deps.Enforcer,
		},
	}))

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	srv.Use(extension.Introspection{})

	// GraphQL playground endpoint
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))

	// GraphQL endpoint
	router.Handle("/query", srv)

	return router
}
