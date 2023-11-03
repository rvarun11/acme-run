package amqp

import "os"

var conf config

type config struct {
	host     string
	port     string
	user     string
	password string
}

func init() {
	conf = newConfig()
}

func newConfig() config {
	return config{
		host:     getEnv("RABBITMQ_HOSTNAME", "localhost"),
		port:     getEnv("RABBITMQ_PORT", "5432"),
		user:     getEnv("RABBITMQ_USER", "guest"),
		password: getEnv("RABBITMQ_PASSWORD", "guest"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
