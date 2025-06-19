package auth

import "github.com/golang-jwt/jwt/v5"

type Authenticator interface {
	GenerateToken(claims *MyClaims) (string, error)
	ValidateToken(token string) (*jwt.Token, *MyClaims, error)
}
