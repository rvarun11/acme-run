package config

import "os"

var Config *AppConfiguration

type AppConfiguration struct {
	Mode     string
	Port     string
	Postgres *Postgres
	RabbitMQ *RabbitMQ
	Publish  *Publish
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	DB_Name  string
	Encoding string
}

type RabbitMQ struct {
	Host     string
	Port     string
	User     string
	Password string
}

type Publish struct {
	Destination string
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

	rabbitmq := &RabbitMQ{
		Host:     getEnv("RABBITMQ_HOSTNAME", "localhost"),
		Port:     getEnv("RABBITMQ_PORT", "5672"),
		User:     getEnv("RABBITMQ_USER", "guest"),
		Password: getEnv("RABBITMQ_PASSWORD", "guest"),
	}

	publish := &Publish{
		Destination: getEnv("SEND_DESTINATION", "SHELTER_TRAIL_WORKOUT"),
	}

	Config = &AppConfiguration{
		Mode:     getEnv("MODE", "dev"),
		Port:     getEnv("PORT", "8005"),
		Postgres: postgres,
		RabbitMQ: rabbitmq,
		Publish:  publish,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
