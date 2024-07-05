package storage

import (
	"github.com/developersismedika/sqlx"

	"github.com/ardihikaru/graphql-example-part-1/pkg/config"
	"github.com/ardihikaru/graphql-example-part-1/pkg/logger"
	mySqlx "github.com/ardihikaru/graphql-example-part-1/pkg/mysqldb"
)

// DbConnect opens MySQL database connection
func DbConnect(log *logger.Logger, dbCfg config.DbMySQL) (*sqlx.DB, error) {
	return mySqlx.DbConnect(log, dbCfg)
}
