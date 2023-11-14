package config

import "os"

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
	DB_Name  string
	Encoding string
}

func init() {
	postgres := &Postgres{
		Host:     getEnv("POSTGRES_HOSTNAME", "localhost"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
		User:     getEnv("POSTGRES_USER", "postgres"),
		Password: getEnv("POSTGRES_PASSWORD", "postgres"),
		DB_Name:  getEnv("POSTGRES_DB", "postgres"),
		Encoding: getEnv("POSTGRES_ENCODING", "UTF8"),
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
