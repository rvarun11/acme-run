package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/CAS735-F23/macrun-teamvsl/challenge/config"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/services"
	logger "github.com/CAS735-F23/macrun-teamvsl/challenge/log"
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

	prefetchCount  = 1
	prefetchSize   = 0
	prefetchGlobal = false

	consumeAutoAck   = true
	consumeExclusive = false
	consumeNoLocal   = false
	consumeNoWait    = false
)

// Images Rabbitmq consumer
type WorkoutStatsConsumer struct {
	amqpConn *amqp.Connection
	svc      *services.ChallengeService
}

func NewWorkoutStatsConsumer(cfg *config.RabbitMQ, challengeSvc *services.ChallengeService) *WorkoutStatsConsumer {

	conn := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	amqpConn, err := amqp.Dial(conn)
	if err != nil {
		logger.Fatal("Unable to dial connection to RabbitMQ")
		return nil
	}

	return &WorkoutStatsConsumer{
		amqpConn: amqpConn,
		svc:      challengeSvc,
	}
}

func (wsc *WorkoutStatsConsumer) InitAMQP() {
	var wg sync.WaitGroup
	wg.Add(1)
	go wsc.StartConsumer(&wg, 1, "", "WORKOUT_STATS_QUEUE", "", "")
}

// Consume messages
func (c *WorkoutStatsConsumer) CreateChannel(exchangeName, queueName, bindingKey, consumerTag string) (*amqp.Channel, error) {
	ch, err := c.amqpConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error amqpConn.Channel %w", err)
	}

	// logger.Debug("Declaring exchange", zap.String("exchange name", exchangeName))
	// err = ch.ExchangeDeclare(
	// 	exchangeName,
	// 	exchangeKind,
	// 	exchangeDurable,
	// 	exchangeAutoDelete,
	// 	exchangeInternal,
	// 	exchangeNoWait,
	// 	nil,
	// )
	// if err != nil {
	// 	return nil, fmt.Errorf("error ch.ExchangeDeclare %w", err)
	// }

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

	logger.Debug("Declaring queue and binding it to exchange",
		zap.String("queue_name", queue.Name),
		zap.String("exchange_name", exchangeName),
		zap.Int("message_count", queue.Messages),
		zap.Int("consumer_count", queue.Consumers),
		zap.String("binding_key", bindingKey),
	)

	// err = ch.QueueBind(
	// 	queue.Name,
	// 	bindingKey,
	// 	exchangeName,
	// 	queueNoWait,
	// 	nil,
	// )
	// if err != nil {
	// 	return nil, fmt.Errorf("error ch.QueueBind %w", err)
	// }

	logger.Debug("Queue bound to exchange, starting to consume from queue", zap.String("consumer_tag", consumerTag))

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
func (c *WorkoutStatsConsumer) StartConsumer(wg *sync.WaitGroup, workerPoolSize int, exchange, queueName, bindingKey, consumerTag string) error {
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
	logger.Debug("ch.NotifyClose", zap.Error(chanErr))
	return chanErr
}

func (c *WorkoutStatsConsumer) worker(ctx context.Context, deliveries <-chan amqp.Delivery) {
	for d := range deliveries {
		csDTO := &challengeStatsDTO{}
		err := json.Unmarshal(d.Body, csDTO)
		// logger.Debug("Received a message: %s", zap.Any("delivery", csDTO))
		if err != nil {
			logger.Debug("failed to unmarshal", zap.Error(err))
		}
		c.svc.CreateOrUpdateChallengeStats(csDTO.PlayerID, csDTO.DistanceCovered, csDTO.EnemiesFought, csDTO.EnemiesEscaped, csDTO.WorkoutEnd)
	}
}
