package main

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/ardihikaru/graphql-example-part-1/internal/application"
	s "github.com/ardihikaru/graphql-example-part-1/internal/storage/user"

	"github.com/ardihikaru/graphql-example-part-1/pkg/config"
	"github.com/ardihikaru/graphql-example-part-1/pkg/logger"
	"github.com/ardihikaru/graphql-example-part-1/pkg/mysqldb"
	"github.com/ardihikaru/graphql-example-part-1/pkg/service/user"
	e "github.com/ardihikaru/graphql-example-part-1/pkg/utils/error"
)

func main() {
	// loads configuration
	cfg, err := config.Get()
	if err != nil {
		e.FatalOnError(err, "failed to load config")
	}

	// builds private key object
	err = cfg.BuildEncryptionKeys()
	if err != nil {
		e.FatalOnError(err, "failed to build private key object")
	}

	// validates config
	err = cfg.Validate()
	if err != nil {
		e.FatalOnError(err, "config validation occurs")
	}

	// configures logger
	log, err := logger.New(cfg.Log.Level, cfg.Log.Format, cfg.Http.LogHttpRequest, &cfg.LogPublisher)
	if err != nil {
		e.FatalOnError(err, "failed to prepare the logger")
	}

	// builds dependencies
	deps := application.BuildDependencies(cfg, log)

	// sample encryption
	plainPasswd := "bismillah"

	// generates hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPasswd), cfg.Encryption.HashCost)
	if err != nil {
		e.FatalOnError(err, "failed to hash the plain password")
	}

	// builds user storage
	userStorage := s.Store{Storage: &mysqldb.Storage{
		Db:  deps.Db,
		Log: deps.Log,
	}}

	// encrypts the plain password
	// FYI: the encrypted plain password will be decrypted, hashed and be compared with the hashed password in the database
	usrSvc := user.NewService(log, &userStorage, cfg)
	encryptedValue, err := usrSvc.EncryptPassword(plainPasswd, cfg.Encryption.PublicKeyRSA)
	if err != nil {
		log.Error(fmt.Sprintf("encryption failed: %s", err.Error()))
	}
	encryptedValueBase64 := base64.StdEncoding.EncodeToString(encryptedValue)

	log.Info(fmt.Sprintf("plainPasswd: %s", plainPasswd))
	log.Info(fmt.Sprintf("hashedPassword: %s", hashedPassword))
	log.Info(fmt.Sprintf("encrypted password base64: %s", encryptedValueBase64))
}
