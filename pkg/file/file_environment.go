package pkgFile

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// path is a route to file .env
func LoadEnv(path string) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Println("No .env file found")
		panic(err)
	}
}

func WriteEnv(path string, key string, value string) error {
	// Append to .env
	envLine := fmt.Sprintf("\n%s=%s\n", key, value)

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
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
