package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TestAuthenticator struct{}

const secret = "test"

var testClaims = jwt.MapClaims{
	"aud": "test-aud",
	"iss": "test-aud",
	"sub": int64(1),
	"exp": time.Now().Add(time.Hour).Unix(),
}

func (a *TestAuthenticator) GenerateToken(claims *MyClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString, nil
}

func (a *TestAuthenticator) ValidateToken(token string) (*jwt.Token, *MyClaims, error) {
	claims := &MyClaims{}
	tok, _ := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	return tok, claims, nil
}
