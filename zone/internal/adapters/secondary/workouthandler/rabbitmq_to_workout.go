package workouthandler

import (
	"encoding/json"
	"fmt"

	"github.com/CAS735-F23/macrun-teamvsl/peripheral/log"
	"github.com/CAS735-F23/macrun-teamvsl/zone/config"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var cfg *config.RabbitMQ = config.Config.RabbitMQ

// declare the name of the queue that publishing data to
var cfgPublish *config.Publish = config.Config.Publish
var destinationQueueName = cfgPublish.Destination

type AMQPPublisher struct {
	amqpConn *amqp.Connection
}

// Initialize new RabbitMQ connection
func NewConnection(cfg *config.RabbitMQ) (*amqp.Connection, error) {
	conn := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)
	mq, err := amqp.Dial(conn)
	if err != nil {
		return &amqp.Connection{}, err
	}

	return mq, nil
}

// NewAMQPPublisher initializes a new AMQPPublisher with a RabbitMQ connection
func NewAMQPPublisher() (*AMQPPublisher, error) {
	log.Debug("creating a trail to workout trail")
	conn, err := NewConnection(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating RabbitMQ connection: %w", err)
	}

	return &AMQPPublisher{amqpConn: conn}, nil
}

// PublishWorkoutStats publishes workout stats to the specified RabbitMQ queue
func (pub *AMQPPublisher) PublishShelterInfo(wId uuid.UUID, sId uuid.UUID, name string, availability bool, distance float64) error {
	ch, err := pub.amqpConn.Channel()
	if err != nil {
		log.Error("publish shelter: failed to open a channel", zap.Error(err))
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()
	var shelter ShelterDTO
	shelter.WorkoutID = wId
	shelter.DistanceToShelter = distance
	shelter.ShelterName = name
	shelter.ShelterID = sId
	shelter.ShelterAvailability = availability
	body, err := json.Marshal(shelter)
	if err != nil {
		log.Error("publish shelter: failed to convert to json data", zap.Error(err))
		return fmt.Errorf("failed to serialize workoutStats: %w", err)
	}
	fmt.Println(destinationQueueName)
	log.Debug("distance to shelter", zap.Float64("distance", distance))
	err = ch.Publish(
		"",                   // exchange
		destinationQueueName, // queue name
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Error("publish shelter: failed to push data", zap.Error(err))
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}
