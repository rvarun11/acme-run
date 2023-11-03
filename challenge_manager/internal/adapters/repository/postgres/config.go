package postgres

import "os"

// Repository Config

var conf config

type config struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
	encoding string
}

func init() {
	conf = newConfig()
}

func newConfig() config {
	return config{
		host:     getEnv("POSTGRES_HOSTNAME", "localhost"),
		port:     getEnv("POSTGRES_PORT", "5432"),
		user:     getEnv("POSTGRES_USER", "postgres"),
		password: getEnv("POSTGRES_PASSWORD", "postgres"),
		dbname:   getEnv("POSTGRES_DB", "postgres"),
		encoding: getEnv("POSTGRES_ENCODING", "UTF8"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
