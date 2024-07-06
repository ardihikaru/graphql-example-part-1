package usergraph

import (
	"github.com/developersismedika/sqlx"

	svc "github.com/ardihikaru/graphql-example-part-1/internal/service/user"
	storage "github.com/ardihikaru/graphql-example-part-1/internal/storage/user"

	"github.com/ardihikaru/graphql-example-part-1/pkg/config"
	"github.com/ardihikaru/graphql-example-part-1/pkg/logger"
	"github.com/ardihikaru/graphql-example-part-1/pkg/mysqldb"
)

func BuildUserService(db *sqlx.DB, log *logger.Logger, cfg *config.Config) *svc.Service {
	return svc.NewService(log, buildUserStorage(db, log), cfg)
}

func buildUserStorage(db *sqlx.DB, log *logger.Logger) *storage.Store {
	// builds user storage
	return &storage.Store{
		Storage: &mysqldb.Storage{
			Db:  db,
			Log: log,
		},
	}
}
