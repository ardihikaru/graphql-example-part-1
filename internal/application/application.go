package application

import (
	"github.com/developersismedika/sqlx"
	"github.com/google/uuid"

	"github.com/ardihikaru/graphql-example-part-1/pkg/config"
	"github.com/ardihikaru/graphql-example-part-1/pkg/enforcer"
	"github.com/ardihikaru/graphql-example-part-1/pkg/jwtauth"
	"github.com/ardihikaru/graphql-example-part-1/pkg/logger"
	"github.com/ardihikaru/graphql-example-part-1/pkg/storage"
	e "github.com/ardihikaru/graphql-example-part-1/pkg/utils/error"
)

// Dependencies hold the primitives and structs and/or interfaces that are required
// for the application's business logic.
type Dependencies struct {
	SvcId     string
	Cfg       *config.Config
	Db        *sqlx.DB
	Log       *logger.Logger
	TokenAuth *jwtauth.JWTAuth
	Enforcer  *enforcer.Enforcer
}

// BuildDependencies builds dependencies
func BuildDependencies(cfg *config.Config, log *logger.Logger) *Dependencies {
	// generates service ID
	svcId := uuid.New().String()

	// initializes JWT Authenticator
	tokenAuth := getTokenAuthentication(&cfg.JwtAuth)

	// initializes persistent store
	db, err := storage.DbConnect(log, cfg.DbMySQL)
	if err != nil {
		e.FatalOnError(err, "failed to open database connection")
	}

	// Load model configuration file and policy store adapter
	enforcerPolicy := enforcer.NewEnforcer(log, db, cfg.Enforcer)
	if err != nil {
		e.FatalOnError(err, "failed to create casbin enforcer")
	}

	return &Dependencies{
		SvcId:     svcId,
		Cfg:       cfg,
		Db:        db,
		Log:       log,
		TokenAuth: tokenAuth,
		Enforcer:  enforcerPolicy,
	}
}
