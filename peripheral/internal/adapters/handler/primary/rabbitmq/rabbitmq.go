package rabbitmqhandler

import (
	"encoding/json"
	"fmt"
	"time"

	// "log"

	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/services"
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/log"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	lastLocationQueueName = "LOCATION_PERIPHERAL_WORKOUT"
	lastHRQueueName       = "Peripheral-HRM-Queue-001"
)

const (
	queueDurable    = true
	queueAutoDelete = false
	queueExclusive  = false
	queueNoWait     = false
)

type RabbitMQHandler struct {
	peripheralService *services.PeripheralService
	amqpURL           string
	connection        *amqp.Connection
	channel           *amqp.Channel
}

func NewRabbitMQHandler(amqpURL string) (*RabbitMQHandler, error) {
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
		amqpURL:    amqpURL,
		connection: conn,
		channel:    ch,
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

func (handler *RabbitMQHandler) SendLastLocation(wId uuid.UUID, latitude float64, longitude float64, time time.Time) error {
	// location := handler.peripheralService.GetGeoLocation(wId)
	var location LastLocation
	location.WorkoutID = wId
	location.Latitude = latitude
	location.Longitude = longitude
	location.TimeOfLocation = time

	body, err := json.Marshal(location)
	if err != nil {
		log.Fatal("Failed to marshal LastLocation:", zap.Error(err))
		return err
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
		return err
	}

	log.Debug("Sent LastLocation", zap.ByteString("body", body))
	return nil
}
