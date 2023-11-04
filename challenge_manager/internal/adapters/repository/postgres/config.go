package postgres

import "github.com/CAS735-F23/macrun-teamvsl/challenge_manager/config"

// Repository Config

var conf *postgresConfig

type postgresConfig struct {
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

func newConfig() *postgresConfig {
	return &postgresConfig{
		host:     config.GetEnv("POSTGRES_HOSTNAME", "localhost"),
		port:     config.GetEnv("POSTGRES_PORT", "5432"),
		user:     config.GetEnv("POSTGRES_USER", "postgres"),
		password: config.GetEnv("POSTGRES_PASSWORD", "postgres"),
		dbname:   config.GetEnv("POSTGRES_DB", "postgres"),
		encoding: config.GetEnv("POSTGRES_ENCODING", "UTF8"),
	}
}
