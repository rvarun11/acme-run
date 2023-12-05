package config

import "os"

var Config *AppConfiguration

type AppConfiguration struct {
	Mode       string
	Port       string
	Postgres   *Postgres
	RabbitMQ   *RabbitMQ
	ZoneClient string
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
	Host                     string
	Port                     string
	User                     string
	Password                 string
	WorkoutLocationPublisher string
	ZoneLocationPublisher    string
}

func init() {
	postgres := &Postgres{
		Host:     getEnv("POSTGRES_HOSTNAME", "localhost"),
		Port:     getEnv("POSTGRES_PORT", "5400"),
		User:     getEnv("POSTGRES_USER", "postgres"),
		Password: getEnv("POSTGRES_PASSWORD", "postgres"),
		DB_Name:  getEnv("POSTGRES_DB", "postgres"),
		Encoding: getEnv("POSTGRES_ENCODING", "UTF8"),
	}

	rabbitmq := &RabbitMQ{
		Host:                     getEnv("RABBITMQ_HOSTNAME", "localhost"),
		Port:                     getEnv("RABBITMQ_PORT", "5672"),
		User:                     getEnv("RABBITMQ_USER", "guest"),
		Password:                 getEnv("RABBITMQ_PASSWORD", "guest"),
		WorkoutLocationPublisher: getEnv("RABBITMQ_WORKOUT_LOCATION_PUBLISHER", "location_peripheral_workout_queue"),
		ZoneLocationPublisher:    getEnv("RABBITMQ_ZONE_LOCATION_PUBLISHER", "location_peripheral_zone_queue"),
	}

	Config = &AppConfiguration{
		Mode:       getEnv("MODE", "prod"),
		Port:       getEnv("PORT", "8012"),
		Postgres:   postgres,
		RabbitMQ:   rabbitmq,
		ZoneClient: getEnv("ZONE_CLIENT_URL", "http://localhost:8005"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
