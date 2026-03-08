package jwt

import "github.com/golang-jwt/jwt/v5"

type AccessClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
