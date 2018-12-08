package web

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/xanderflood/fruit-pi/lib/tools"
)

type Authorization = jwt.Claims

//Authenticator Authenticator
//go:generate counterfeiter . Authenticator
type Authenticator interface {
	Authenticate(w http.ResponseWriter, r *http.Request, claimsObj Authorization) bool
}

//AuthenticatorAgent AuthenticatorAgent
type AuthenticatorAgent struct {
	Logger        tools.Logger
	SigningSecret []byte
}

//NewAuthenticator NewAuthenticator
func NewAuthenticator(logger tools.Logger, signingSecret string) *AuthenticatorAgent {
	return &AuthenticatorAgent{
		Logger:        logger,
		SigningSecret: []byte(signingSecret),
	}
}

//Authenticate checks and unmarshals the bearer token
func (a *AuthenticatorAgent) Authenticate(w http.ResponseWriter, r *http.Request, claimsObj Authorization) bool {
	var tokenString string
	authHeader := r.Header.Get("Authorization")

	_, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenString)
	if err != nil {
		a.Logger.Error(err)
		AuthorizationError(w, "invalid `Authorization` header: expected bearer token")
		return false
	}

	_, err = jwt.ParseWithClaims(tokenString, claimsObj, a.keyFunc)
	if err != nil {
		AuthorizationError(w, "invalid token")
		return false
	}

	return true
}

//AuthorizationError AuthorizationError
func AuthorizationError(w http.ResponseWriter, message string) {
	w.Header().Set("WWW-Authenticate", message)
	JSONErrorRespond(w, message, http.StatusUnauthorized)
}

func (a *AuthenticatorAgent) keyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(a.SigningSecret), nil
}
