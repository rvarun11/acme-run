package amqphandler

import (
	"encoding/json"
	"fmt"

	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/services"
	log "github.com/CAS735-F23/macrun-teamvsl/workout/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		failMsg := fmt.Sprintf("%s: %s", msg, err)
		log.Error(failMsg)
	}
}

func ShelterSubscriber(svc services.WorkoutService, url string) {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"SHELTER_QUEUE_001", // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	failOnError(err, "Failed to declare SHELTER_QUEUE_001")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register ShelterSubscriber")

	var forever chan struct{}

	var shelterAvailable ShelterAvailable

	go func() {
		for d := range msgs {
			log.Info(fmt.Sprintf("Received a message: %s", d.Body))
			err = json.Unmarshal(d.Body, &shelterAvailable)
			svc.UpdateShelter(shelterAvailable.WorkoutID, shelterAvailable.DistanceToShelter)
		}
	}()

	log.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func LocationSubscriber(svc services.WorkoutService, url string) {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"LOCATION_QUEUE_001", // name
		false,                // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	failOnError(err, "Failed to declare LOCATION_QUEUE_001")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register LocationSubscriber")

	var forever chan struct{}

	var lastLocation LastLocation

	go func() {
		for d := range msgs {
			log.Info(fmt.Sprintf("Received a message: %s", d.Body))
			err = json.Unmarshal(d.Body, &lastLocation)
			svc.UpdateDistanceTravelled(lastLocation.WorkoutID, lastLocation.Latitude, lastLocation.Longitude, lastLocation.TimeOfLocation)
		}
	}()

	log.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
