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

// NewPublisher initializes a new Publisher with a RabbitMQ connection
func NewPublisher(cfg *config.RabbitMQ) *Publisher {
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

	return &Publisher{
		amqpConn: amqpConn,
	}
}

// PublishWorkoutStats publishes workout stats to the specified RabbitMQ queue
func (pub *Publisher) PublishWorkoutStats(workoutStats *domain.Workout) error {
	ch, err := pub.amqpConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	// Declare the queue to ensure it exists
	_, err = ch.QueueDeclare(
		"WORKOUT_STATS_QUEUE", // queue name
		true,                  // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

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
