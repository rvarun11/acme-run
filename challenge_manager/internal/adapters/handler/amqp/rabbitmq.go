package amqp

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	mq *amqp.Connection
}

// Initialize new RabbitMQ connection
func NewConnection() (*RabbitMQ, error) {
	conn := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		conf.user,
		conf.password,
		conf.host,
		conf.port,
	)
	mq, err := amqp.Dial(conn)
	if err != nil {
		return &RabbitMQ{}, err
	}

	return &RabbitMQ{
		mq: mq,
	}, nil
}

// Images Rabbitmq consumer
type WorkoutStatsConsumer struct {
	logger log.Logger
	// emailUC  email.EmailsUseCase
}
