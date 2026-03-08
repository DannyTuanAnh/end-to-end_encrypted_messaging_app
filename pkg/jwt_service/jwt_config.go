package jwt

import "time"

type JWTConfig struct {
	Issuer string

	AccessSecret  string
	RefreshSecret string

	AccessTokenExpire  time.Duration
	RefreshTokenExpire time.Duration
}
