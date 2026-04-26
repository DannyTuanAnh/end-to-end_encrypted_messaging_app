package jwt_app

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignJWT(privateKeyPEM string, claims jwt.MapClaims) (string, error) {
	keyData, err := os.ReadFile(privateKeyPEM)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return "", fmt.Errorf("invalid pem")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return "", err
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signed, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func CreatePayload(issService string, userID int64, audService string) jwt.MapClaims {
	return jwt.MapClaims{
		"iss":     issService,
		"user_id": userID,
		"aud":     audService,
		"exp":     time.Now().Add(5 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}
}
