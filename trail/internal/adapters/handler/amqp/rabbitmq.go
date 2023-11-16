package amqp

import (
	"encoding/json"
	"log"

	"github.com/CAS735-F23/macrun-teamvsl/trail/internal/core/ports"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	repoTM  ports.TrailManagerRepository
}

func NewRabbitMQService(amqpURL string, repoTM ports.TrailManagerRepository) (*RabbitMQService, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &RabbitMQService{
		conn:    conn,
		channel: channel,
		repoTM:  repoTM,
	}, nil
}

func (r *RabbitMQService) ListenForLocationUpdates(queueName string, wId uuid.UUID) {
	msgs, err := r.channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var location LocationDTO
			if err := json.Unmarshal(d.Body, &location); err != nil {
				log.Printf("Error decoding JSON: %s", err)
				continue
			}

			// Assume UpdateLocation is a method that updates the TrailManager's location
			if err := r.repoTM.UpdateLocation(location.WorkoutID, location.Longitude, location.Latitude); err != nil {
				log.Printf("Failed to update location: %s", err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (r *RabbitMQService) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

// package amqphandler

// import (
// 	"context"
// 	"encoding/json"
// 	"log"
// 	"math"
// 	"time"

// 	"github.com/google/uuid"
// 	amqp "github.com/rabbitmq/amqp091-go"
// )

// func failOnError(err error, msg string) {
// 	if err != nil {
// 		log.Panicf("%s: %s", msg, err)
// 	}
// }

// func PeripheralSubscriber(svc services.TrailManagerService, url string) {
// 	conn, err := amqp.Dial(url)
// 	failOnError(err, "Failed to connect to RabbitMQ")
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	failOnError(err, "Failed to open a channel")
// 	defer ch.Close()

// 	q, err := ch.QueueDeclare(
// 		"HR-Queue-001", // name
// 		false,          // durable
// 		false,          // delete when unused
// 		false,          // exclusive
// 		false,          // no-wait
// 		nil,            // arguments
// 	)
// 	failOnError(err, "Failed to declare a queue")

// 	msgs, err := ch.Consume(
// 		q.Name, // queue
// 		"",     // consumer
// 		true,   // auto-ack
// 		false,  // exclusive
// 		false,  // no-local
// 		false,  // no-wait
// 		nil,    // args
// 	)
// 	failOnError(err, "Failed to register a consumer")

// 	var forever chan struct{}

// 	var tempDTOVar LastLocation
// 	go func() {
// 		for d := range msgs {
// 			log.Printf("Received a message: %s", d.Body)
// 			err = json.Unmarshal(d.Body, &tempDTOVar)
// 			// TODO: Ignoring Error for now, Handle Error later
// 			// Call the following to get the HR Value updated
// 			svc.CurrentLatitude = tempDTOVar.Latitude
// 			svc.CurrentLongitude = tempDTOVar.Longitude
// 			svc.CurrentTime = TimeOfLocation
// 			//svc.UpdateHRValue(tempDTOVar.WorkoutID, tempDTOVar.HRValue)

// 		}
// 	}()

// 	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
// 	<-forever
// }

// func publishShelter(workoutId uuid.UUID, shelterId uuid.UUID, distance float64) {

// 	// Just connect for now and send
// 	// TODO: Should we connect once and use the same for sending and receiving?

// 	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5903/")
// 	failOnError(err, "Failed to connect to RabbitMQ")
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	failOnError(err, "Failed to open a channel")
// 	defer ch.Close()

// 	q, err := ch.QueueDeclare(
// 		"TRAIL-Workout-001", // name
// 		false,               // durable
// 		false,               // delete when unused
// 		false,               // exclusive
// 		false,               // no-wait
// 		nil,                 // arguments
// 	)
// 	failOnError(err, "Failed to declare a queue")

// 	failOnError(err, "Failed to declare a queue")
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	var tempDTOVar ShelterAvailable
// 	tempDTOVar.WorkoutID = workoutId
// 	tempDTOVar.ShelterAvailable = ((shelterId == uuid.nil) || (distance == math.MaxFloat64))
// 	tempDTOVar.DistanceToShelter = distance

// 	var body []byte

// 	body, _ = json.Marshal(tempDTOVar)

// 	err = ch.PublishWithContext(ctx,
// 		"",     // exchange
// 		q.Name, // routing key
// 		false,  // mandatory
// 		false,  // immediate
// 		amqp.Publishing{
// 			ContentType: "text/json",
// 			Body:        []byte(body),
// 		})
// 	failOnError(err, "Failed to publish shelter vailable")
// 	log.Printf(" [x] Sent %s\n", body)
// }
