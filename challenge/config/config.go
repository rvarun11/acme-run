package config

import "os"

var Config *AppConfiguration

type AppConfiguration struct {
	Mode     string
	Port     string
	Postgres *Postgres
	RabbitMQ *RabbitMQ
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

type RabbitMQ struct {
	Host                 string
	Port                 string
	User                 string
	Password             string
	WorkoutStatsConsumer string
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

	rabbitmq := &RabbitMQ{
		Host:                 getEnv("RABBITMQ_HOSTNAME", "localhost"),
		Port:                 getEnv("RABBITMQ_PORT", "5672"),
		User:                 getEnv("RABBITMQ_USER", "guest"),
		Password:             getEnv("RABBITMQ_PASSWORD", "guest"),
		WorkoutStatsConsumer: getEnv("RABBITMQ_WORKOUT_STATS_CONSUMER", "WORKOUT_STATS_QUEUE"),
	}

	Config = &AppConfiguration{
		Mode:     getEnv("MODE", "dev"),
		Port:     getEnv("PORT", "8003"),
		Postgres: postgres,
		RabbitMQ: rabbitmq,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
