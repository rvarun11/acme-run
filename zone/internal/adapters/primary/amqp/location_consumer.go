package amqp

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/CAS735-F23/macrun-teamvsl/zone/config"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/core/services"
	logger "github.com/CAS735-F23/macrun-teamvsl/zone/log"
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

	consumeAutoAck   = false
	consumeExclusive = false
	consumeNoLocal   = false
	consumeNoWait    = false
)

// Location AMQP consumer
type LocationConsumer struct {
	amqpConn *amqp.Connection
	svc      *services.ZoneService
	config   *config.RabbitMQ
}

func NewLocationConsumer(cfg *config.RabbitMQ, zoneSvc *services.ZoneService) *LocationConsumer {
	conn := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	amqpConn, err := amqp.Dial(conn)
	if err != nil {
		logger.Fatal("unable to dial connection to RabbitMQ")
	}

	return &LocationConsumer{
		amqpConn: amqpConn,
		svc:      zoneSvc,
		config:   cfg,
	}
}

func (lc *LocationConsumer) InitAMQP() {
	var wg sync.WaitGroup
	wg.Add(1)
	go lc.StartConsumer(&wg, 1, "", lc.config.LiveLocationConsumer, "", "")
}

func (lc *LocationConsumer) CreateChannel(exchangeName, queueName, bindingKey, consumerTag string) (*amqp.Channel, error) {
	ch, err := lc.amqpConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error amqpConn.Channel %w", err)
	}

	// logger.Debug("declaring exchange", zap.String("exchange name", exchangeName))
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

func (lc *LocationConsumer) StartConsumer(wg *sync.WaitGroup, workerPoolSize int, exchange, queueName, bindingKey, consumerTag string) error {
	// Remove the context cancellation, as this should only be done when the application stops
	ch, err := lc.CreateChannel(exchange, queueName, bindingKey, consumerTag)
	if err != nil {
		return fmt.Errorf("create channel error %w", err)
	}

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
		ch.Close()
		return fmt.Errorf("consume error %w", err)
	}

	for i := 0; i < workerPoolSize; i++ {
		logger.Debug("Starting worker", zap.Int("worker number", i))
		go lc.worker(deliveries)
	}

	// Do not close the channel here, it will be closed when the application exits
	//defer ch.Close()
	return nil
}

func (lc *LocationConsumer) worker(deliveries <-chan amqp.Delivery) {
	for d := range deliveries {
		var lastLocation LocationDTO
		err := json.Unmarshal(d.Body, &lastLocation)
		if err != nil {
			logger.Error("Failed to unmarshal", zap.Error(err))
			d.Nack(false, false) // negatively acknowledge the message and requeue it if needed
			continue
		}

		logger.Debug("Received a message and unmarshalled successfully", zap.Any("location", lastLocation))
		// Process the message...
		err = lc.svc.UpdateCurrentLocation(lastLocation.WorkoutID, lastLocation.Latitude, lastLocation.Longitude, lastLocation.TimeOfLocation)
		if err != nil {
			logger.Error("Failed to update current location", zap.Error(err))
		}
		d.Ack(false) // acknowledge the message upon successful processing
	}
}
