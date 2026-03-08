package jwt

import (
	"context"
	"errors"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Validator interface {
	ValidateAccess(ctx context.Context, tokenString string) (*AccessClaims, error)
	ValidateRefresh(ctx context.Context, tokenString string) (*jwtlib.RegisteredClaims, error)
}

type hs256Validator struct {
	cfg JWTConfig
}

func NewValidator(cfg JWTConfig) Validator {
	return &hs256Validator{cfg: cfg}
}

func (v *hs256Validator) ValidateAccess(ctx context.Context, tokenString string) (*AccessClaims, error) {

	claims := &AccessClaims{}

	token, err := jwtlib.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwtlib.Token) (interface{}, error) {
			return []byte(v.cfg.AccessSecret), nil
		},
		jwtlib.WithIssuer(v.cfg.Issuer), // check iss
	)

	if err != nil {
		// Token hết hạn
		if errors.Is(err, jwtlib.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}

		//Issuer sai
		if errors.Is(err, jwtlib.ErrTokenInvalidIssuer) {
			return nil, ErrInvalidIssuer
		}

		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (v *hs256Validator) ValidateRefresh(ctx context.Context, tokenString string) (*jwtlib.RegisteredClaims, error) {

	claims := &jwtlib.RegisteredClaims{}

	token, err := jwtlib.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwtlib.Token) (interface{}, error) {
			return []byte(v.cfg.RefreshSecret), nil
		},
		jwtlib.WithIssuer(v.cfg.Issuer), // check iss
	)

	if err != nil {
		// Token hết hạn
		if errors.Is(err, jwtlib.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}

		//Issuer sai
		if errors.Is(err, jwtlib.ErrTokenInvalidIssuer) {
			return nil, ErrInvalidIssuer
		}

		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
