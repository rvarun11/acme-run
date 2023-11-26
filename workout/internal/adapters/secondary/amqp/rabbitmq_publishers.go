package amqphandler

import (
	"encoding/json"
	"fmt"

	"github.com/CAS735-F23/macrun-teamvsl/workout/config"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/domain"
	logger "github.com/CAS735-F23/macrun-teamvsl/workout/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Publisher struct {
	amqpConn *amqp.Connection
}

var cfg *config.RabbitMQ = config.Config.RabbitMQ

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

// NewPublisher initializes a new Publisher with a RabbitMQ connection
func NewPublisher() (*Publisher, error) {
	conn, err := NewConnection(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating RabbitMQ connection: %w", err)
	}

	return &Publisher{amqpConn: conn}, nil
}

// PublishWorkoutStats publishes workout stats to the specified RabbitMQ queue
func (pub *Publisher) PublishWorkoutStats(workoutStats *domain.Workout) error {
	ch, err := pub.amqpConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	var challengeStatsDTO = challengeStatsDTO{
		PlayerID:        workoutStats.PlayerID,
		WorkoutEnd:      workoutStats.EndedAt,
		EnemiesFought:   workoutStats.Fights,
		EnemiesEscaped:  workoutStats.Escapes,
		DistanceCovered: workoutStats.DistanceCovered,
	}

	body, err := json.Marshal(challengeStatsDTO)
	if err != nil {
		return fmt.Errorf("failed to serialize workoutStats: %w", err)
	}

	err = ch.Publish(
		"",                    // exchange
		"WORKOUT_STATS_QUEUE", // queue name
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	logger.Info("Workout statistics published to Challenge Manager", zap.Any("Stats", challengeStatsDTO))
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}
