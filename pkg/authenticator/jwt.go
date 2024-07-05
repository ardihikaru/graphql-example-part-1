// Package authenticator provides functions to authenticate incoming requests
package authenticator

import (
	"github.com/lestrrat-go/jwx/jwa"
	"golang.org/x/crypto/bcrypt"

	"github.com/ardihikaru/graphql-example-part-1/pkg/jwtauth"
)

// JWTClaims is the data type of claims in Encode() function (jwtauth.go)
type JWTClaims map[string]interface{}

// MakeTokenAuth creates token authentication
func MakeTokenAuth(jwtAlgo jwa.SignatureAlgorithm, jwtSecret string) (*jwtauth.JWTAuth, error) {
	// generates token auth
	tokenAuth := jwtauth.New(jwtAlgo.String(), []byte(jwtSecret), nil)

	return tokenAuth, nil
}

// CreateAccessToken creates an access token
func CreateAccessToken(tokenAuth *jwtauth.JWTAuth, jwtClaims JWTClaims) string {
	// extracts token
	_, tokenString, _ := tokenAuth.Encode(jwtClaims)

	return tokenString
}

// CheckPasswordHash compares the input password with the stored password,
// and finally validates if they are match or not
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
