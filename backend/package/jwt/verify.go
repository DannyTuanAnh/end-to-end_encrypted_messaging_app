package jwt_app

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyJWT(certPath string, tokenString string) (*jwt.Token, error) {
	var block *pem.Block

	is_cloud_run := os.Getenv("IS_CLOUD_RUN")

	if is_cloud_run == "true" {
		block, _ = pem.Decode([]byte(certPath))
	} else {

		certData, err := os.ReadFile(certPath)
		if err != nil {
			return nil, err
		}

		block, _ = pem.Decode(certData)
	}

	if block == nil {
		return nil, fmt.Errorf("invalid cert")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey := cert.PublicKey

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodRS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return publicKey, nil
	})

	return token, err
}
