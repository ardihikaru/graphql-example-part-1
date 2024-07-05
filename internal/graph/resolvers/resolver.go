package graph

import (
	"sync"

	"github.com/developersismedika/sqlx"

	"github.com/ardihikaru/graphql-example-part-1/pkg/config"
	"github.com/ardihikaru/graphql-example-part-1/pkg/enforcer"
	"github.com/ardihikaru/graphql-example-part-1/pkg/jwtauth"
	"github.com/ardihikaru/graphql-example-part-1/pkg/logger"
)

type Resolver struct {
	// All messages since launching the GraphQL endpoint

	// All active subscriptions
	mu sync.Mutex

	// supporting variables
	SvcId     string
	Cfg       *config.Config
	Db        *sqlx.DB
	Log       *logger.Logger
	TokenAuth *jwtauth.JWTAuth
	Enforcer  *enforcer.Enforcer
}
