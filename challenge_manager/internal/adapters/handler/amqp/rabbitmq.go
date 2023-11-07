package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/CAS735-F23/macrun-teamvsl/challenge_manager/config"
	"github.com/CAS735-F23/macrun-teamvsl/challenge_manager/internal/core/services"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangeKind       = "direct"
	exchangeDurable    = true
	exchangeAutoDelete = false
	exchangeInternal   = false
	exchangeNoWait     = false

	queueDurable    = true
	queueAutoDelete = false
	queueExclusive  = false
	queueNoWait     = false

	// publishMandatory = false
	// publishImmediate = false

	prefetchCount  = 1
	prefetchSize   = 0
	prefetchGlobal = false

	consumeAutoAck   = false
	consumeExclusive = false
	consumeNoLocal   = false
	consumeNoWait    = false
)

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

// Images Rabbitmq consumer
type WorkoutStatsConsumer struct {
	amqpConn *amqp.Connection
	logger   log.Logger
	svc      services.ChallengeService
}

// Consume messages
func (c *WorkoutStatsConsumer) CreateChannel(exchangeName, queueName, bindingKey, consumerTag string) (*amqp.Channel, error) {
	ch, err := c.amqpConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error amqpConn.Channel %w", err)
	}

	c.logger.Printf("Declaring exchange: %s", exchangeName)
	err = ch.ExchangeDeclare(
		exchangeName,
		exchangeKind,
		exchangeDurable,
		exchangeAutoDelete,
		exchangeInternal,
		exchangeNoWait,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error ch.ExchangeDeclare %w", err)
	}

	queue, err := ch.QueueDeclare(
		queueName,
		queueDurable,
		queueAutoDelete,
		queueExclusive,
		queueNoWait,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("error ch.QueueDeclare %w", err)
	}

	c.logger.Printf("Declared queue, binding it to exchange: Queue: %v, messagesCount: %v, "+
		"consumerCount: %v, exchange: %v, bindingKey: %v",
		queue.Name,
		queue.Messages,
		queue.Consumers,
		exchangeName,
		bindingKey,
	)

	err = ch.QueueBind(
		queue.Name,
		bindingKey,
		exchangeName,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error ch.QueueBind %w", err)
	}

	c.logger.Printf("Queue bound to exchange, starting to consume from queue, consumerTag: %v", consumerTag)

	err = ch.Qos(
		prefetchCount,  // prefetch count
		prefetchSize,   // prefetch size
		prefetchGlobal, // global
	)
	if err != nil {
		return nil, fmt.Errorf("error ch.Qos %w", err)
	}

	return ch, nil
}

// Start new rabbitmq consumer
func (c *WorkoutStatsConsumer) StartConsumer(workerPoolSize int, exchange, queueName, bindingKey, consumerTag string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := c.CreateChannel(exchange, queueName, bindingKey, consumerTag)
	if err != nil {
		return fmt.Errorf("create channel error %w", err)
	}
	defer ch.Close()

	deliveries, err := ch.Consume(
		queueName,
		consumerTag,
		consumeAutoAck,
		consumeExclusive,
		consumeNoLocal,
		consumeNoWait,
		nil,
	)
	if err != nil {
		return fmt.Errorf("consume error %w", err)
	}

	for i := 0; i < workerPoolSize; i++ {
		/// Do something with the deliveriesFind
		go c.worker(ctx, deliveries)
	}

	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	c.logger.Fatalf("ch.NotifyClose: %v", chanErr)
	return chanErr
}

func (c *WorkoutStatsConsumer) worker(ctx context.Context, deliveries <-chan amqp.Delivery) {
	for d := range deliveries {
		var csDTO *challengeStatsDTO
		log.Printf("Received a message: %s", d.Body)
		err := json.Unmarshal(d.Body, csDTO)
		if err != nil {
			log.Panicf("failed to unmarshal %s", err)
		}
		c.svc.SubscribeToActiveChallenges(csDTO.PlayerID, csDTO.DistanceCovered, csDTO.EnemiesFought, csDTO.EnemiesEscaped)
	}
}
