package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/nyaruka/phonenumbers"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

var (
	pathEnv = ".env"
)

// loadEnv is a helper function that loads environment variables from a .env file using the godotenv package
func LoadEnv() {
	err := godotenv.Load(pathEnv)
	if err != nil {
		log.Println("No .env file found")
		panic(err)
	}
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func GetEnvTime(key string, defaultValue int) time.Duration {
	if value := os.Getenv(key); value != "" {
		val, err := strconv.Atoi(value)
		if err != nil {
			log.Println("Error converting environment variable to int:", err)
		} else {
			return time.Duration(val)
		}
	}

	return time.Duration(defaultValue)
}

func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		val, err := strconv.Atoi(value)
		if err != nil {
			log.Println("Error converting environment variable to int:", err)
		} else {
			return val
		}
	}

	return defaultValue
}

func GetEnvList(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		valTrim := strings.TrimSpace(value)
		if valTrim == "" {
			return defaultValue
		}
		return strings.Split(valTrim, ",")
	}

	return defaultValue
}

func SetListBoolean(blocked []string) map[string]bool {
	blockedDomains := make(map[string]bool)
	for _, domain := range blocked {
		blockedDomains[NormalizeString(domain)] = true
	}

	return blockedDomains
}

func WriteEnv(key string, value string) error {
	// Append to .env
	envLine := fmt.Sprintf("\n%s=%s\n", key, value)

	f, err := os.OpenFile(pathEnv, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(envLine); err != nil {
		return err
	}
	log.Printf("%s=%s", key, value)

	return nil
}

func SaveKeyToEnv(clientType, key string) error {
	if clientType == "web" {
		return WriteEnv("VITE_API_KEY", key)
	}

	envKey := fmt.Sprintf("%s_API_KEY", strings.ToUpper(clientType))

	log.Println("Saving API key to environment variable:", envKey)

	return WriteEnv(envKey, key)
}

func CheckUUID(id string) bool {
	err := uuid.Validate(id)
	if err != nil {
		return false
	}

	idUuid, err := uuid.Parse(id)
	if err != nil {
		return false
	}

	if idUuid == uuid.Nil {
		return false
	}

	return true
}

// grpc, mtls, get caller service name from client certificate
func GetCaller(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok || p == nil {
		return "unknown"
	}

	tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return "unknown"
	}

	if len(tlsInfo.State.PeerCertificates) == 0 {
		return "unknown"
	}

	cert := tlsInfo.State.PeerCertificates[0]

	if len(cert.DNSNames) > 0 {
		return cert.DNSNames[0]
	}

	return "unknown"
}

func IsPhoneNumber(phone *string) (string, bool) {
	num, err := phonenumbers.Parse(*phone, "VN")
	if err != nil {
		return "", false
	}

	if !phonenumbers.IsValidNumber(num) || phonenumbers.GetNumberType(num) != phonenumbers.MOBILE {
		return "", false
	}

	e164 := phonenumbers.Format(num, phonenumbers.E164)
	return e164, true
}
