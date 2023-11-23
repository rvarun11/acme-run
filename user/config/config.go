package config

import (
	"os"
)

var Config *AppConfiguration

type AppConfiguration struct {
	Mode     string
	Port     string
	Postgres *Postgres
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Encoding string
	LogLevel string
}

func init() {
	postgres := &Postgres{
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
		User:     getEnv("POSTGRES_USER", "postgres"),
		Password: getEnv("POSTGRES_PASSWORD", "postgres"),
		DBName:   getEnv("POSTGRES_DB", "postgres"),
		Encoding: getEnv("POSTGRES_ENCODING", "UTF8"),
		LogLevel: getEnv("POSTGRES_LOG_LEVEL", "warn"),
	}

	Config = &AppConfiguration{
		Mode:     getEnv("MODE", "dev"),
		Port:     getEnv("PORT", "8000"),
		Postgres: postgres,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
