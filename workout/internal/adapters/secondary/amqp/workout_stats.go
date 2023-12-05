package amqp

import (
	"encoding/json"
	"fmt"

	"github.com/CAS735-F23/macrun-teamvsl/workout/config"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/domain"
	logger "github.com/CAS735-F23/macrun-teamvsl/workout/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	// exchangeKind       = "direct"
	// exchangeDurable    = true
	// exchangeAutoDelete = false
	// exchangeInternal   = false
	// exchangeNoWait     = false

	queueDurable    = true
	queueAutoDelete = false
	queueExclusive  = false
	queueNoWait     = false

	// publishMandatory = false
	// publishImmediate = false

	// prefetchCount  = 1
	// prefetchSize   = 0
	// prefetchGlobal = false

	// consumeAutoAck   = true
	// consumeExclusive = false
	// consumeNoLocal   = false
	// consumeNoWait    = false
)

type WorkoutStatsPublisher struct {
	amqpConn *amqp.Connection
	config   *config.RabbitMQ
}

// NewWorkoutStatsPublisher initializes a new WorkoutStatsPublisher with a RabbitMQ connection
func NewWorkoutStatsPublisher(cfg *config.RabbitMQ) *WorkoutStatsPublisher {
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

	return &WorkoutStatsPublisher{
		config:   cfg,
		amqpConn: amqpConn,
	}
}

// PublishWorkoutStats publishes workout stats to the specified RabbitMQ queue
func (pub *WorkoutStatsPublisher) PublishWorkoutStats(workoutStats *domain.Workout) error {
	ch, err := pub.amqpConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	// Declare the queue to ensure it exists
	_, err = ch.QueueDeclare(
		pub.config.WorkoutStatsPublisher, // queue name
		queueDurable,                     // durable
		queueAutoDelete,                  // delete when unused
		queueExclusive,                   // exclusive
		queueNoWait,                      // no-wait
		nil,                              // arguments
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
		"",                               // exchange
		pub.config.WorkoutStatsPublisher, // queue name
		false,                            // mandatory
		false,                            // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	logger.Info("workout statistics published", zap.Any("stats", challengeStatsDTO))
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}
