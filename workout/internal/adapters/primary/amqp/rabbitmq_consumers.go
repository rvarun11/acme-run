package amqphandler

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

var cfg *config.RabbitMQ = config.Config.RabbitMQ

// Shelter Rabbitmq consumer
type ShelterSubscriber struct {
	amqpConn *amqp.Connection
	svc      *services.WorkoutService
}

// Location Rabbitmq consumer
type LocationSubscriber struct {
	amqpConn *amqp.Connection
	svc      *services.WorkoutService
}

type WorkoutAMQPHandler struct {
	locationSubscriber *LocationSubscriber
	shelterSubscriber  *ShelterSubscriber
}

func NewAMQPHandler(workoutSvc *services.WorkoutService) *WorkoutAMQPHandler {
	amqpConn_shelterSub, err := NewConnection(cfg)
	if err != nil {
		logger.Error("Connection to RabbitMQ Failed")
	}
	shelterSubscriber := ShelterSubscriber{
		amqpConn: amqpConn_shelterSub,
		svc:      workoutSvc,
	}

	amqpConn_locationSub, err := NewConnection(cfg)
	if err != nil {
		logger.Error("Connection to RabbitMQ Failed")
	}
	locationSubscriber := LocationSubscriber{
		amqpConn: amqpConn_locationSub,
		svc:      workoutSvc,
	}

	return &WorkoutAMQPHandler{
		locationSubscriber: &locationSubscriber,
		shelterSubscriber:  &shelterSubscriber,
	}
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

// Consume messages
func (c *ShelterSubscriber) CreateChannel(exchangeName, queueName, bindingKey, consumerTag string) (*amqp.Channel, error) {
	ch, err := c.amqpConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error amqpConn.Channel %w", err)
	}

	/*logger.Debug("Declaring exchange", zap.String("exchange name", exchangeName))
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
	}*/

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
func (c *ShelterSubscriber) StartConsumer(wg *sync.WaitGroup, workerPoolSize int, exchange, queueName, bindingKey, consumerTag string) error {
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

func (c *ShelterSubscriber) worker(ctx context.Context, deliveries <-chan amqp.Delivery) {
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

// Consume messages
func (c *LocationSubscriber) CreateChannel(exchangeName, queueName, bindingKey, consumerTag string) (*amqp.Channel, error) {
	ch, err := c.amqpConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error amqpConn.Channel %w", err)
	}

	logger.Debug("Declaring exchange", zap.String("exchange name", exchangeName))
	/*err = ch.ExchangeDeclare(
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
	}*/

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
func (c *LocationSubscriber) StartConsumer(wg *sync.WaitGroup, workerPoolSize int, exchange, queueName, bindingKey, consumerTag string) error {

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

	// TODO Fix blocking error handling
	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	logger.Debug("ch.NotifyClose", zap.Error(chanErr))
	return nil
}

func (c *LocationSubscriber) worker(ctx context.Context, deliveries <-chan amqp.Delivery) {
	for d := range deliveries {
		lastLocation := &LastLocation{}
		logger.Debug("Received a message: ", zap.String("delivery", string(d.Body)))
		err := json.Unmarshal(d.Body, lastLocation)
		if err != nil {
			logger.Debug("failed to unmarshal %s", zap.Error(err))
		}
		c.svc.UpdateDistanceTravelled(lastLocation.WorkoutID, lastLocation.Latitude, lastLocation.Longitude, lastLocation.TimeOfLocation)
	}
}

func (wah *WorkoutAMQPHandler) InitAMQP() error {
	var wg sync.WaitGroup
	wg.Add(2)
	go wah.locationSubscriber.StartConsumer(&wg, 1, "", "LOCATION_PERIPHERAL_WORKOUT", "", "")
	go wah.shelterSubscriber.StartConsumer(&wg, 1, "", "SHELTER_TRAIL_WORKOUT", "", "")

	return nil
}
