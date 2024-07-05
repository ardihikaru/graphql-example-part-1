package enforcer

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	sqlxadapter "github.com/developersismedika/casbin-sqlx-adapter"
	"github.com/developersismedika/sqlx"

	"github.com/ardihikaru/graphql-example-part-1/pkg/config"
	"github.com/ardihikaru/graphql-example-part-1/pkg/logger"
	e "github.com/ardihikaru/graphql-example-part-1/pkg/utils/error"
)

type Enforcer struct {
	*casbin.Enforcer
	log *logger.Logger
}

// NewEnforcer creates a new Enforcer object
//
//	other possible adapter can be found here: https://casbin.org/docs/adapters/
func NewEnforcer(log *logger.Logger, db *sqlx.DB, enforcerCfg config.Enforcer) *Enforcer {
	// builds adapter options
	opts := &sqlxadapter.AdapterOptions{
		TableName: enforcerCfg.TableName,

		// reuses an existing connection:
		DB: db,
	}
	adapter := sqlxadapter.NewAdapterFromOptions(opts)

	// loads model configuration file and policy store adapter
	enforcer, err := casbin.NewEnforcer(enforcerCfg.ModelFile, adapter)
	if err != nil {
		e.FatalOnError(err, fmt.Sprintf("failed to create casbin enforcer"))
	}

	return &Enforcer{
		Enforcer: enforcer,
		log:      log,
	}
}
