package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

// Claims struct for JWT
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
