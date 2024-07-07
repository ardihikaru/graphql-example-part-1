package auth

import (
	"time"

	"github.com/ardihikaru/graphql-example-part-1/pkg/authenticator"
	"github.com/ardihikaru/graphql-example-part-1/pkg/jwtauth"
	"github.com/ardihikaru/graphql-example-part-1/pkg/logger"
	"github.com/ardihikaru/graphql-example-part-1/pkg/service/session"
)

const (
	ClaimExpiredInKey = "exp"
	ClaimIssuedAtKey  = "iat"
	ClaimUserKey      = "user"
)

type Token struct {
	AccessToken string          `json:"access_token"`
	ExpiredIn   int64           `json:"expired_in"`
	IssuedAt    int64           `json:"issued_at"`
	Session     session.Session `json:"session"`
}

// LoginData is the input JSON body captured from the login request
type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Service prepares the interfaces related with this auth service
type Service struct {
	log        *logger.Logger
	tokenAuth  *jwtauth.JWTAuth
	jwtExpTime int
}

// NewService creates a new auth service
func NewService(log *logger.Logger, jwtExpTime int, tokenAuth *jwtauth.JWTAuth) *Service {
	return &Service{
		log:        log,
		jwtExpTime: jwtExpTime,
		tokenAuth:  tokenAuth,
	}
}

// Authorize authorizes user credential
func (svc *Service) Authorize(userId int, username string) (*Token, error) {
	claimUser := session.ClaimUser{
		UserId:    int64(userId),
		AccountId: "id123",
		Username:  username,
		Name:      "Ini Budi",
		Role:      "administrator", // dummy value
		Email:     "user@email.com",
	}

	// builds the JWT claim options
	durationIn := time.Duration(svc.jwtExpTime) * time.Second
	jwtClaims := authenticator.JWTClaims{
		ClaimExpiredInKey: jwtauth.ExpireIn(durationIn),
		ClaimIssuedAtKey:  jwtauth.EpochNow(),
		ClaimUserKey:      claimUser,

		// TODO: adds other claims
	}

	// begins to create the access token
	accessToken := authenticator.CreateAccessToken(svc.tokenAuth, jwtClaims)

	token := &Token{
		AccessToken: accessToken,
		ExpiredIn:   jwtauth.ExpireIn(durationIn),
		IssuedAt:    jwtauth.EpochNow(),
		Session: session.Session{
			AccountId: "id123",
			UserEmail: "user@email.com",
			UserId:    int64(userId),
			Username:  username,
			Role:      "administrator", // dummy value
			Name:      "Ini Budi",
		},
	}

	return token, nil
}
