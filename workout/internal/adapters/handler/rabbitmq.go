package handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/CAS735-F23/macrun-teamvs_/workout/internal/core/services"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func HRMReceiver(svc services.WorkoutService) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"HR-Queue-001", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	type tempDTO struct {
		WorkoutID uuid.UUID `json:"workoutID"`
		HRValue   uint16    `json:"hrValue"`
	}

	var tempDTOVar tempDTO
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			err = json.Unmarshal(d.Body, &tempDTOVar)
			// TODO: Ignoring Error for now, Handle Error later
			// Call the following to get the HR Value updated
			svc.UpdateHRValue(tempDTOVar.WorkoutID, tempDTOVar.HRValue)

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func TieWorkoutToHRM(hrmID uuid.UUID, workoutID uuid.UUID) {

	// Just connect for now and send
	// TODO: Should we connect once and use the same for sending and receiving?

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"HR-Workout-001", // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	failOnError(err, "Failed to declare a queue")

	failOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: Temp DTO
	type tempDTO struct {
		WorkoutID uuid.UUID `json:"workoutID"`
		HRMId     uuid.UUID `json:"hrmID"`
	}

	var tempDTOVar tempDTO
	tempDTOVar.WorkoutID = workoutID
	tempDTOVar.HRMId = hrmID

	var body []byte

	body, err = json.Marshal(tempDTOVar)

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)

}
