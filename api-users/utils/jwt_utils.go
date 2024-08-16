package utils

type JWTUtils interface {
	GenerateToken(userID string) (string, error)
}
