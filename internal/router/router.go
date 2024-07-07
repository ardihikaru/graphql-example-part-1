package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	"github.com/ardihikaru/graphql-example-part-1/internal/application"
	"github.com/ardihikaru/graphql-example-part-1/internal/router/graphql"

	"github.com/ardihikaru/graphql-example-part-1/pkg/logger"
)

// GetRouter configures a chi router and starts the http server
func GetRouter(deps *application.Dependencies) *chi.Mux {
	r := chi.NewRouter()

	if deps.Log != nil {
		r.Use(logger.SetLogger(deps.Log))
	}

	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   deps.Cfg.Cors.AllowedOrigins,
		AllowedMethods:   deps.Cfg.Cors.AllowedMethods,
		AllowedHeaders:   deps.Cfg.Cors.AllowedHeaders,
		ExposedHeaders:   deps.Cfg.Cors.ExposedHeaders,
		AllowCredentials: deps.Cfg.Cors.AllowCredentials,
		MaxAge:           deps.Cfg.Cors.MaxAge,
		Debug:            deps.Cfg.Cors.Debug,
	}))

	buildGraphQLTree(r, deps)

	return r
}

// buildGraphQLTree handles graphQL related route(s)
func buildGraphQLTree(r *chi.Mux, deps *application.Dependencies) {
	r.Mount("/", gqlrouter.GraphQLRoutes(deps))
}
