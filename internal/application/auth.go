package application

import (
	"github.com/ardihikaru/graphql-example-part-1/pkg/authenticator"
	"github.com/ardihikaru/graphql-example-part-1/pkg/config"
	"github.com/ardihikaru/graphql-example-part-1/pkg/jwtauth"
	e "github.com/ardihikaru/graphql-example-part-1/pkg/utils/error"
)

// getTokenAuthentication creates an authentication token from the authenticator
func getTokenAuthentication(jwtCfg *config.JwtAuth) *jwtauth.JWTAuth {
	tokenAuth, err := authenticator.MakeTokenAuth(jwtCfg.Algorithm, jwtCfg.Secret)
	if err != nil {
		e.FatalOnError(err, "failed to create a JWT authenticator")
	}

	return tokenAuth
}
