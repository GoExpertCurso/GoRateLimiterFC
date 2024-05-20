package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	requiredEnv := []string{"DB_DRIVER", "DB_HOST", "DB_PORT", "DB_PASSWORD", "IP_RATE_LIMIT", "IP_BLOCK_DURATION", "TOKEN_RATE_LIMIT", "TOKEN_BLOCK_DURATION"}
	for _, env := range requiredEnv {
		if os.Getenv(env) == "" {
			log.Fatalf("Environment variable %s not set", env)
		}
	}
}
