package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	HTTP_PORT            string
	PRODUCT_SERVICE_PORT string
	AUTH_SERVICE_PORT    string
	DB_HOST              string
	DB_PORT              string
	DB_USER              string
	DB_PASSWORD          string
	DB_NAME              string
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	config := &Config{}

	config.PRODUCT_SERVICE_PORT = cast.ToString(coalesce("PRODUCT_SERVICE_PORT", ":8070"))
	config.DB_HOST = cast.ToString(coalesce("DB_HOST", "aritsanConnP"))
	config.DB_PORT = cast.ToString(coalesce("DB_PORT", 5432))
	config.DB_USER = cast.ToString(coalesce("DB_USER", "postgres"))
	config.DB_PASSWORD = cast.ToString(coalesce("DB_PASSWORD", "pass"))
	config.DB_NAME = cast.ToString(coalesce("DB_NAME", "artisanconnect_proto"))

	return config
}

func coalesce(key string, value interface{}) interface{} {
	val, exists := os.LookupEnv(key)
	if exists {
		return val
	}
	return value
}
