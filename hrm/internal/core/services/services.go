package services

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/CAS735-F23/macrun-teamvsl/hrm/internal/core/ports"
	"github.com/google/uuid"

	"github.com/CAS735-F23/macrun-teamvsl/hrm/internal/core/domain"
)



type HRMService struct {
	repo ports.HRMRepository
}


func NewHRMService(repo ports.HRMRepository) *HRMService {
	return &HRMService{
		repo: repo,
	}
}

func (s *HRMService) ConnectHRM(hrmID uuid.UUID) {
	var h domain.HRM
	//var err error
	h.HRMId = hrmID
	hrm, _ := domain.NewHRM(h)
	s.repo.AddHRMIntance(hrm)
}

func (s *HRMService) DisconnectHRM(hrmID uuid.UUID) {
	s.repo.DeleteHRMInstance(hrmID)
}

func (s *HRMService) BindHRMtoWorkout(hrmID uuid.UUID, workoutID uuid.UUID) {
	//var err error
	hrmInstance, _ := s.repo.Get(hrmID)
	//TODO Error Handling
	hrmInstance.WorkoutId = workoutID
	s.repo.Update(*hrmInstance)
}

func (s *HRMService) SendHRM(wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq/")
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

	failOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: Temp DTO
	type tempDTO struct {
		WorkoutID uuid.UUID `json:"workoutID"`
		HRValue   int       `json:"hrValue"`
	}
	for {
		hrms, _ := s.repo.List()
		for i := 0; i < len(hrms); i++ {
			min := 30
			max := 200
			var tempDTOVar tempDTO
			tempDTOVar.WorkoutID = (*hrms[i]).WorkoutId
			tempDTOVar.HRValue = rand.Intn(max-min) + min

			var body []byte

			body, _ = json.Marshal(tempDTOVar)

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
		time.Sleep(5 * time.Second)
	}

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
