package handler

import (
	"encoding/json"
	"fmt"

	// "log"

	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/dto"
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/services"
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	lastLocationQueueName = "LOCATION_PERIPHERAL_WORKOUT"
	lastHRQueueName       = "Peripheral-HRM-Queue-001"
)

// LS-TODO: Remove warnings

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

	publishMandatory = false
	publishImmediate = false

	prefetchCount  = 1
	prefetchSize   = 0
	prefetchGlobal = false

	consumeAutoAck   = false
	consumeExclusive = false
	consumeNoLocal   = false
	consumeNoWait    = false
)

// LS-TODO: Please look at challenge service Rabbit MQ handler and see how it's being done.
type RabbitMQHandler struct {
	peripheralService *services.PeripheralService
	amqpURL           string
	connection        *amqp.Connection
	channel           *amqp.Channel
}

func NewRabbitMQHandler(p *services.PeripheralService, amqpURL string) (*RabbitMQHandler, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	handler := &RabbitMQHandler{
		peripheralService: p,
		amqpURL:           amqpURL,
		connection:        conn,
		channel:           ch,
	}

	// Ensure the queues exist when starting up
	err = handler.declareQueues()
	if err != nil {
		handler.Close()
		return nil, err
	}

	return handler, nil
}

func (handler *RabbitMQHandler) declareQueues() error {
	queues := []string{lastLocationQueueName, lastHRQueueName}
	for _, queue := range queues {
		_, err := handler.channel.QueueDeclare(
			queue,
			queueDurable,
			queueAutoDelete,
			queueExclusive,
			queueNoWait,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue '%s': %w", queue, err)
		}
	}
	return nil
}

func (handler *RabbitMQHandler) Close() {
	if handler.channel != nil {
		handler.channel.Close()
	}
	if handler.connection != nil {
		handler.connection.Close()
	}
}

func (handler *RabbitMQHandler) SendLastLocation(tLoc dto.LastLocation) {
	// location := handler.peripheralService.GetGeoLocation(wId)
	body, err := json.Marshal(tLoc)
	if err != nil {
		log.Fatal("Failed to marshal LastLocation:", zap.Error(err))
		return
	}

	err = handler.channel.Publish(
		"",
		lastLocationQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Fatal("Failed to publish LastLocation", zap.Error(err))
		return
	}

	log.Debug("Sent LastLocation", zap.ByteString("body", body))
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal(msg, zap.Error(err))
	}
}

// package handler

// import (
// 	"encoding/json"
// 	"log"

// 	// "github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/ports"

// 	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/services"
// 	"github.com/google/uuid"
// 	amqp "github.com/rabbitmq/amqp091-go"
// )

// const (
// 	exchangeKind       = "direct"
// 	exchangeDurable    = true
// 	exchangeAutoDelete = false
// 	exchangeInternal   = false
// 	exchangeNoWait     = false

// 	queueDurable    = true
// 	queueAutoDelete = false
// 	queueExclusive  = false
// 	queueNoWait     = false

// 	// publishMandatory = false
// 	// publishImmediate = false

// 	prefetchCount  = 1
// 	prefetchSize   = 0
// 	prefetchGlobal = false

// 	consumeAutoAck   = false
// 	consumeExclusive = false
// 	consumeNoLocal   = false
// 	consumeNoWait    = false
// )

// const (
// 	lastLocationQueueName = "LOCATION_PERIPHERAL_WORKOUT"
// 	lastHRQueueName       = "Peripheral-HRM-Queue-001"
// )

// type RabbitMQHandler struct {
// 	peripheralService *services.PeripheralService
// 	amqpURL           string
// }

// func NewRabbitMQHandler(p *services.PeripheralService, amqpURL string) *RabbitMQHandler {
// 	return &RabbitMQHandler{
// 		peripheralService: p,
// 		amqpURL:           amqpURL,
// 	}
// }

// func (handler *RabbitMQHandler) SendLastLocation(wId uuid.UUID) {
// 	conn, err := amqp.Dial(handler.amqpURL)
// 	failOnError(err, "Failed to connect to RabbitMQ")
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	failOnError(err, "Failed to open a channel")
// 	defer ch.Close()

// 	q, err := ch.QueueDeclare(
// 		lastLocationQueueName, // name
// 		false,                 // durable
// 		false,                 // delete when unused
// 		false,                 // exclusive
// 		false,                 // no-wait
// 		nil,                   // arguments
// 	)
// 	failOnError(err, "Failed to declare a queue")

// 	location := handler.peripheralService.GetGeoLocation(wId)
// 	body, err := json.Marshal(location)
// 	failOnError(err, "Failed to marshal LastLocation")

// 	err = ch.Publish(
// 		"",     // exchange
// 		q.Name, // routing key
// 		false,  // mandatory
// 		false,  // immediate
// 		amqp.Publishing{
// 			ContentType: "application/json",
// 			Body:        body,
// 		})
// 	failOnError(err, "Failed to publish LastLocation")

// 	log.Printf("Sent LastLocation: %s\n", body)
// }

// func failOnError(err error, msg string) {
// 	if err != nil {
// 		log.Fatalf("%s: %s", msg, err)
// 	}
// }
