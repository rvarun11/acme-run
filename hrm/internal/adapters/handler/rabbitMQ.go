package handler

import (
	"encoding/json"
	"log"

	"github.com/CAS735-F23/macrun-teamvsl/hrm/internal/core/services"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func HRMWorkoutBinder(svc services.HRMService) {
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

	// TODO: Temp DTO
	type tempDTO struct {
		WorkoutID uuid.UUID `json:"workoutID"`
		HRMId     uuid.UUID `json:"hrmID"`
	}

	var tempDTOVar tempDTO
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			err = json.Unmarshal(d.Body, &tempDTOVar)
			failOnError(err, "Failed to unmarshal")
			// TODO: Ignoring Error for now, Handle Error later
			// Call the following to get the HR Value updated
			svc.BindHRMtoWorkout(tempDTOVar.HRMId, tempDTOVar.WorkoutID)

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
