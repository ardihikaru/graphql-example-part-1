package usergraph

import (
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"

	"github.com/ardihikaru/graphql-example-part-1/internal/application"
	"github.com/ardihikaru/graphql-example-part-1/internal/service/auth"
	s "github.com/ardihikaru/graphql-example-part-1/internal/storage/user"

	"github.com/ardihikaru/graphql-example-part-1/pkg/graph/user/generated"
	graph "github.com/ardihikaru/graphql-example-part-1/pkg/graph/user/resolvers"
	"github.com/ardihikaru/graphql-example-part-1/pkg/mysqldb"
	"github.com/ardihikaru/graphql-example-part-1/pkg/service/user"
)

const (
	BaseRoute = "/user"
)

// BuildRouter builds the router
func BuildRouter(deps *application.Dependencies) http.Handler {
	r := chi.NewRouter()

	buildGraphQLHandler(r, deps)

	return r
}

// buildGraphQLHandler builds graphQL handler
func buildGraphQLHandler(r *chi.Mux, deps *application.Dependencies) {
	// builds user storage
	userStorage := s.Store{Storage: &mysqldb.Storage{
		Db:  deps.Db,
		Log: deps.Log,
	}}

	// registers services
	authSvc := auth.NewService(deps.Log, deps.Cfg.JwtAuth.ExpiredInSec, deps.TokenAuth)
	userSvc := user.NewService(deps.Log, &userStorage, deps.Cfg)

	// creates GraphQL server
	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: graph.NewResolver(deps.Cfg, deps.Db, deps.Log, deps.TokenAuth, deps.Enforcer,
			userSvc, authSvc),
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

	// registers GraphQL playground endpoint
	r.Handle("/", playground.Handler("GraphQL playground", fmt.Sprintf("%s/query", BaseRoute)))

	// registers GraphQL endpoint
	r.Handle("/query", srv)
}
