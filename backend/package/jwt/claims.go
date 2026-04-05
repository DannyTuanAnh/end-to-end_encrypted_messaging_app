package jwt_app

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}
