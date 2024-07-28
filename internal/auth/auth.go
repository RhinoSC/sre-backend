// auth/auth.go
package auth

import (
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

func Init(secret string) {
	tokenAuth = jwtauth.New("HS256", []byte(secret), nil)
}

func GenerateToken(userID string) (string, error) {
	claims := map[string]interface{}{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expiration time
	}
	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Verifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(tokenAuth)
}

func Authenticator() func(http.Handler) http.Handler {
	return jwtauth.Authenticator(tokenAuth)
}
