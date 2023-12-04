package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/CAS735-F23/macrun-teamvsl/workout/config"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/services"
	logger "github.com/CAS735-F23/macrun-teamvsl/workout/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// Shelter AMQP Consumer
type ShelterDistanceConsumer struct {
	amqpConn *amqp.Connection
	svc      *services.WorkoutService
	config   *config.RabbitMQ
}

func NewShelterDistanceConsumer(cfg *config.RabbitMQ, workoutSvc *services.WorkoutService) *ShelterDistanceConsumer {
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

	return &ShelterDistanceConsumer{
		config:   cfg,
		amqpConn: amqpConn,
		svc:      workoutSvc,
	}
}

func (wsc *ShelterDistanceConsumer) InitAMQP() {
	var wg sync.WaitGroup
	wg.Add(1)
	go wsc.StartConsumer(&wg, 1, "", wsc.config.ShelterDistanceConsumer, "", "")
}

// Consume messages
func (c *ShelterDistanceConsumer) CreateChannel(exchangeName, queueName, bindingKey, consumerTag string) (*amqp.Channel, error) {
	ch, err := c.amqpConn.Channel()
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

	/*err = ch.QueueBind(
		queue.Name,
		bindingKey,
		exchangeName,
		queueNoWait,
		nil,
	)*/
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

// Start new rabbitmq consumer
func (c *ShelterDistanceConsumer) StartConsumer(wg *sync.WaitGroup, workerPoolSize int, exchange, queueName, bindingKey, consumerTag string) error {
	defer wg.Done()
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

	var forever chan struct{}
	for i := 0; i < workerPoolSize; i++ {
		/// Do something with the deliveriesFind
		go c.worker(ctx, deliveries)
	}

	// TODO Fix blocking error handling
	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	logger.Debug("ch.NotifyClose", zap.Error(chanErr))
	<-forever
	return nil
}

func (c *ShelterDistanceConsumer) worker(ctx context.Context, deliveries <-chan amqp.Delivery) {
	for d := range deliveries {
		shelterAvailable := &ShelterAvailable{}
		logger.Debug("Received a message: ", zap.String("delivery", string(d.Body)))
		err := json.Unmarshal(d.Body, shelterAvailable)
		if err != nil {
			logger.Debug("failed to unmarshal %s", zap.Error(err))
		}
		c.svc.UpdateShelter(shelterAvailable.WorkoutID, shelterAvailable.DistanceToShelter)
	}
}
