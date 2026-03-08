package jwt

import "errors"

var (
	ErrMissingToken  = errors.New("missing token")
	ErrInvalidToken  = errors.New("invalid token")
	ErrExpiredToken  = errors.New("expired token")
	ErrInvalidIssuer = errors.New("invalid issuer")
)
