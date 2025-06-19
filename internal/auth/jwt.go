package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// MyClaims defines the custom JWT claims structure for type safety
type MyClaims struct {
	jwt.RegisteredClaims
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type JWTAuthenticator struct {
	secret string
	aud    string
	iss    string
}

func NewJWTAuthenticator(secret, aud, iss string) *JWTAuthenticator {
	if secret == "" {
		panic("JWT secret must not be empty")
	}
	return &JWTAuthenticator{
		secret: secret,
		aud:    aud,
		iss:    iss,
	}
}

func (a *JWTAuthenticator) GenerateToken(claims *MyClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *JWTAuthenticator) ValidateToken(tokenString string) (*jwt.Token, *MyClaims, error) {
	claims := &MyClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(a.secret), nil
	}, jwt.WithExpirationRequired(),
		jwt.WithAudience(a.aud),
		jwt.WithIssuer(a.iss),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
	if err != nil {
		return nil, nil, err
	}
	if !token.Valid {
		return nil, nil, fmt.Errorf("invalid token")
	}
	return token, claims, nil
}
