package auth

import "github.com/golang-jwt/jwt/v5"

type JWTValidator interface {
	ValidateToken(tokenString string, secret []byte) (jwt.MapClaims, error)
}
