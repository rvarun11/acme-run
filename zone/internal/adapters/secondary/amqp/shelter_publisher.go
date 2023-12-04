package amqp

import (
	"encoding/json"
	"fmt"

	"github.com/CAS735-F23/macrun-teamvsl/zone/config"
	logger "github.com/CAS735-F23/macrun-teamvsl/zone/log"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// declare the name of the queue that publishing data to
type ShelterDistancePublisher struct {
	amqpConn *amqp.Connection
	config   *config.RabbitMQ
}

// NewShelterDistancePublisher initializes a new ShelterDistancePublisher with a RabbitMQ connection
func NewShelterDistancePublisher(cfg *config.RabbitMQ) *ShelterDistancePublisher {
	conn := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	amqpConn, err := amqp.Dial(conn)
	if err != nil {
		logger.Fatal("unable to dial connection to RabbitMQ", zap.Error(err))
		return nil
	}

	return &ShelterDistancePublisher{
		config:   cfg,
		amqpConn: amqpConn,
	}
}

// PublishWorkoutStats publishes workout stats to the specified RabbitMQ queue
func (pub *ShelterDistancePublisher) PublishShelterDistance(wId uuid.UUID, sId uuid.UUID, name string, availability bool, distance float64) error {
	ch, err := pub.amqpConn.Channel()
	if err != nil {
		logger.Error("publish shelter: failed to open a channel", zap.Error(err))
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
		logger.Error("publish shelter: failed to convert to json data", zap.Error(err))
		return fmt.Errorf("failed to serialize workoutStats: %w", err)
	}
	logger.Debug("distance to shelter", zap.Float64("distance", distance))
	err = ch.Publish(
		"",                                  // exchange
		pub.config.ShelterDistancePublisher, // queue name
		false,                               // mandatory
		false,                               // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		logger.Error("publish shelter: failed to push data", zap.Error(err))
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}
