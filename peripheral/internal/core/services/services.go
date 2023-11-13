package services

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/ports"
	"github.com/google/uuid"

	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/domain"
)

type PeripheralService struct {
	repo ports.PeripheralRepository
}

func NewPeripheralService(repo ports.PeripheralRepository) *PeripheralService {
	return &PeripheralService{
		repo: repo,
	}
}

func (s *PeripheralService) ConnectPeripheral(wId uuid.UUID, hrmId uuid.UUID) {
	var h domain.Peripheral
	//var err error
	h.HRMId = hrmID
	h.WorkoutId = wId
	p, _ := domain.NewPeripheral(h)
	s.repo.AddPeripheralIntance(p)
}

func (s *PeripheralService) DisconnectPeripheral(wId uuid.UUID) {
	s.repo.DeletePeripheralInstance(wId)
}

func (s *PeripheralService) SendPeripheralId(wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"Peripheral-Queue-001", // name
		false,                  // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
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
		ps, _ := s.repo.List()
		for i := 0; i < len(ps); i++ {
			min := 30
			max := 200
			var tempDTOVar tempDTO
			tempDTOVar.WorkoutID = (*ps[i]).WorkoutId
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

func (s *PeripheralService) GetHRMReading(wId uuid.UUID) LastHR {
	pInstance, _ := s.repo.Get(wId)
	return pInstance.GetAverageHRate()
}

func (s *PeripheralService) SetHRMReading(wId uuid.UUID, rate int) {
	pInstance, _ := s.repo.Get(wId)
	return pInstance.SetAverageHRate(rate)
}

func (s *PeripheralService) GetHRMDevStatus(wId uuid.UUID) bool {
	pInstance, _ := s.repo.Get(wId)
	return pInstance.GetHRMStatus()
}

func (s *PeripheralService) SetHRMDevStatus(wId uuid.UUID, code bool) {
	pInstance, _ := s.repo.Get(wId)
	pInstance.SetHRMStatus(code)
}

func (s *PeripheralService) SetGeoLocation(wId uuid.UUID, longitude float64, latitude float64) {
	pInstance, _ := s.repo.Get(wId)
	pInstance.SetLocation(longitude, latitude)
}

func (s *PeripheralService) GetGeoDevStatus(wId uuid.UUID) bool {
	pInstance, _ := s.repo.Get(wId)
	return pInstance.GetGeoStatus()
}

func (s *PeripheralService) SetGeoDevStatus(wId uuid.UUID, code bool) {
	pInstance, _ := s.repo.Get(wId)
	pInstance.SetGeoStatus(code)
}

// func (s *PeripheralService) SendPeripheralLocation(*sync.WaitGroup) {
// 	defer wg.Done()
// 	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq/")
// 	failOnError(err, "Failed to connect to RabbitMQ")
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	failOnError(err, "Failed to open a channel")
// 	defer ch.Close()

// 	q, err := ch.QueueDeclare(
// 		"Peripheral-Queue-002", // name
// 		false,                  // durable
// 		false,                  // delete when unused
// 		false,                  // exclusive
// 		false,                  // no-wait
// 		nil,                    // arguments
// 	)
// 	failOnError(err, "Failed to declare a queue")

// 	failOnError(err, "Failed to declare a queue")
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	for {
// 		ps, _ := s.repo.List()
// 		for i := 0; i < len(ps); i++ {
// 			min := 30
// 			max := 200
// 			tempDTOVar := (*ps[i]).getLocation()
// 			var body []byte

// 			body, _ = json.Marshal(tempDTOVar)

// 			err = ch.PublishWithContext(ctx,
// 				"",     // exchange
// 				q.Name, // routing key
// 				false,  // mandatory
// 				false,  // immediate
// 				amqp.Publishing{
// 					ContentType: "text/json",
// 					Body:        []byte(body),
// 				})
// 			failOnError(err, "Failed to publish a message")
// 			log.Printf(" [x] Sent %s\n", body)
// 		}
// 		time.Sleep(5 * time.Second)
// 	}
// }

func (s *PeripheralService) GetGeoLocation(wId uuid.UUID) {
	pInstance, _ := s.repo.Get(wId)
	return pInstance.getLocation()
}
