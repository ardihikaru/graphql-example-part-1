package user

import (
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/ardihikaru/graphql-example-part-1/internal/graph/model"
	"github.com/ardihikaru/graphql-example-part-1/internal/service/user/dto"

	"github.com/ardihikaru/graphql-example-part-1/pkg/config"
	"github.com/ardihikaru/graphql-example-part-1/pkg/logger"
)

// storage provides the interface for the functionality of MinIO
type storage interface {
	InsertUser(username, hashedPassword, status string, isAdmin int, setCreateUserID int64) (*dto.User, error)
	GetUserById(userId int64) (*dto.User, error)
	GetUsers(userIdStr, statusCd string) ([]*dto.User, error)
	GetUserCredByUsername(username string) (*dto.User, *string, error)
}

// Service load user service with model.User
type Service struct {
	log     *logger.Logger
	storage storage
	cfg     *config.Config
}

// NewService creates new user service
func NewService(log *logger.Logger, storage storage, cfg *config.Config) *Service {
	return &Service{
		log:     log,
		storage: storage,
		cfg:     cfg,
	}
}

// Create creates a new user
func (svc *Service) Create(data model.UserInput, CreateUserID *int64) (*model.User, error) {
	var err error
	var ret model.User

	userId, err := svc.GetUserIdByUsername(data.UserNm)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if userId != 0 {
		return nil, fmt.Errorf("username is exist")
	}

	var setCreateUserID int64
	if CreateUserID != nil {
		setCreateUserID = *CreateUserID
	}

	var passDecrypted []byte
	passDecrypted, err = svc.decryptPassword(data.PassHash)
	if err != nil {
		return &ret, err
	}
	hashedPassword, _ := hashPassword(string(passDecrypted), svc.cfg.Encryption.HashCost)

	userRecord, err := svc.storage.InsertUser(data.UserNm, hashedPassword, "active", data.IsAdmin,
		setCreateUserID)

	if err != nil {
		return &ret, err
	}

	return userRecord.ToModel(), nil
}

// List gets list of user
func (svc *Service) List(userIdStr, statusCd string) ([]*model.User, error) {
	userList, err := svc.storage.GetUsers(userIdStr, statusCd)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil, nil when no rows are found
		}
		return nil, err // Return other errors
	}

	var userListModel []*model.User
	for _, usr := range userList {
		userListModel = append(userListModel, usr.ToModel())
	}

	return userListModel, nil
}

// GetById gets user data by ID
func (svc *Service) GetById(userId int64) (*model.User, error) {
	usr, err := svc.storage.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	return usr.ToModel(), nil
}

// GetUserIdByUsername gets user data by username
func (svc *Service) GetUserIdByUsername(username string) (int, error) {
	var err error

	userData, _, err := svc.storage.GetUserCredByUsername(username)
	if err != nil {
		return 0, err
	}

	userId, err := strconv.Atoi(userData.UserID)

	return userId, nil
}

// Authenticate authenticates the provided credential
func (svc *Service) Authenticate(userName string, password string) (bool, error) {
	start := time.Now()
	var err error
	var hashedPassword *string

	_, hashedPassword, err = svc.storage.GetUserCredByUsername(userName)
	if err != nil {
		return false, err
	}
	svc.log.Debug(fmt.Sprintf("hashedPassword: %s", *hashedPassword))

	if err != nil {
		return false, err
	}

	svc.log.Debug(fmt.Sprintf("Execution took %s", time.Since(start)))

	var passDecrypted []byte
	passDecrypted, err = svc.decryptPassword(password)
	if err != nil {
		svc.log.Error(fmt.Sprintf("decrypt: %s", err.Error()))
		return false, fmt.Errorf("DECRYPTION_FAILED")
	}

	if svc.checkPasswordHash(*hashedPassword, string(passDecrypted)) {
		return true, nil
	} else {
		return false, fmt.Errorf("PASSWORD_MISMATCH")
	}
}

// hashPassword hashes given password
func hashPassword(password string, hashCost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		return "", fmt.Errorf("HASH_PASSWORD_ERROR")
	}
	return string(bytes), nil
}

// EncryptPassword encrypts password
func (svc *Service) EncryptPassword(password string, publicKeyRSA *rsa.PublicKey) ([]byte, error) {
	passEncrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKeyRSA, []byte(password))
	if err != nil {
		svc.log.Error(fmt.Sprintf("encryption failed: %s", err.Error()))
		return nil, fmt.Errorf("ENCRYPTION_FAILED")
	}

	return passEncrypted, nil
}

// decryptPassword decrypts the encoded password
//
//	openssl genrsa -out private.pem 1024
//	openssl rsa -in private.pem -outform PEM -pubout -out public.pem
func (svc *Service) decryptPassword(passEncoded string) ([]byte, error) {
	var err error

	var passDecrypted []byte
	var passDecoded []byte
	passDecoded, err = base64.StdEncoding.DecodeString(passEncoded)
	if err != nil {
		log.Printf("base64 decode: %s\n", err)
		return nil, fmt.Errorf("BASE64_DECODE_FAILED")
	}

	passDecrypted, err = rsa.DecryptPKCS1v15(rand.Reader, svc.cfg.Encryption.PrivateKeyRSA, []byte(passDecoded))
	if err != nil {
		log.Printf("decrypt: %s\n", err)
		return nil, err
	}

	return passDecrypted, nil
}

// checkPasswordHash compares password hash
func (svc *Service) checkPasswordHash(hash, password string) bool {
	start := time.Now()

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	svc.log.Debug(fmt.Sprintf("CompareHashAndPassword execution took %s", time.Since(start)))
	return err == nil
}
