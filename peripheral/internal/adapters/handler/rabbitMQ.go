package amqphandler

import (
	"encoding/json"
	"log"

	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/ports"
	"github.com/google/uuid"
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

const (
	lastLocationQueueName = "Peripheral-Location-Queue-001"
	lastHRQueueName       = "Peripheral-HRM-Queue-001"
)

type RabbitMQHandler struct {
	peripheralService *ports.PeripheralService
	amqpURL           string
}

func NewRabbitMQHandler(peripheralService *ports.PeripheralService, amqpURL string) *RabbitMQHandler {
	return &RabbitMQHandler{
		peripheralService: peripheralService,
		amqpURL:           amqpURL,
	}
}

func (handler *RabbitMQHandler) SendLastLocation(wId uuid.UUID) {
	conn, err := amqp.Dial(handler.amqpURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		lastLocationQueueName, // name
		false,                 // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	location := handler.peripheralService.GetGeoLocation(wId)
	body, err := json.Marshal(location)
	failOnError(err, "Failed to marshal LastLocation")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish LastLocation")

	log.Printf("Sent LastLocation: %s\n", body)
}

func (handler *RabbitMQHandler) SendLastHR(wId uuid.UUID) {
	conn, err := amqp.Dial(handler.amqpURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		lastHRQueueName, // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	hrmReading := handler.peripheralService.GetHRMReading(wId)
	body, err := json.Marshal(hrmReading)
	failOnError(err, "Failed to marshal LastHR")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish LastHR")

	log.Printf("Sent LastHR: %s\n", body)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
