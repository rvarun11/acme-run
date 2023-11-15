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

func randomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func (s *PeripheralService) CreatePeripheral(pId uuid.UUID, hId uuid.UUID) {
	var h domain.Peripheral
	h.HRMId = hId
	h.PlayerId = pId
	p, _ := domain.NewPeripheral(h)
	s.repo.AddPeripheralIntance(p)
}

func (s *PeripheralService) CheckStatusByHRMId(hId uuid.UUID) bool {
	pInstance, err := s.repo.GetByHRMId(hId)
	if err != nil {
		return false
	} else {
		pInstance.HRMId = hId
		s.repo.Update(*pInstance)
		return true
	}
}

func (s *PeripheralService) BindPeripheral(pId uuid.UUID, wId uuid.UUID, hrmId uuid.UUID, connected bool, sendToTrail bool) {

	pInstance, err := s.repo.GetByHRMId(hrmId)
	if err != nil {
		var h domain.Peripheral
		h.HRMId = hrmId
		h.WorkoutId = wId
		h.LiveData = true
		h.GeoBrodacasting = sendToTrail
		h.HRMStatus = connected
		h.PlayerId = pId
		p, _ := domain.NewPeripheral(h)
		s.repo.AddPeripheralIntance(p)
		return
	}

	pInstance.PlayerId = pId
	pInstance.WorkoutId = wId
	pInstance.SetHRMStatus(connected)
	pInstance.LiveData = true
	pInstance.GeoBrodacasting = sendToTrail
	s.repo.Update(*pInstance)

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

func (s *PeripheralService) GetHRMAvgReading(hId uuid.UUID) domain.LastHR {
	pInstance, _ := s.repo.GetByHRMId(hId)
	return pInstance.GetAverageHRate()
}

func (s *PeripheralService) GetHRMReading(hId uuid.UUID) domain.LastHR {
	pInstance, _ := s.repo.GetByHRMId(hId)
	return pInstance.GetHRate()
}

func (s *PeripheralService) SetHeartRateReading(hId uuid.UUID, reading int) {
	pInstance, error := s.repo.GetByHRMId(hId)
	if error != nil {
		log.Fatal("error updating HRM %v", error)

		return
	}
	pInstance.SetHRate(reading)
	s.repo.Update(*pInstance)
}

func (s *PeripheralService) GetHRMDevStatus(wId uuid.UUID) bool {
	pInstance, _ := s.repo.Get(wId)
	return pInstance.GetHRMStatus()
}

func (s *PeripheralService) GetHRMDevStatusByHRMId(hId uuid.UUID) bool {
	pInstance, _ := s.repo.GetByHRMId(hId)
	return pInstance.GetHRMStatus()
}

func (s *PeripheralService) SetHRMDevStatusByHRMId(hId uuid.UUID, code bool) {
	pInstance, _ := s.repo.GetByHRMId(hId)
	pInstance.SetHRMStatus(code)
	s.repo.Update(*pInstance)
}

func (s *PeripheralService) SetHRMDevStatus(wId uuid.UUID, code bool) {
	pInstance, _ := s.repo.Get(wId)
	pInstance.SetHRMStatus(code)
	s.repo.Update(*pInstance)
}

func (s *PeripheralService) SetGeoLocation(wId uuid.UUID, longitude float64, latitude float64) {
	pInstance, _ := s.repo.Get(wId)
	pInstance.SetLocation(longitude, latitude)
	s.repo.Update(*pInstance)
}

func (s *PeripheralService) GetGeoDevStatus(wId uuid.UUID) bool {
	pInstance, _ := s.repo.Get(wId)
	return pInstance.GetGeoStatus()
}

func (s *PeripheralService) SetGeoDevStatus(wId uuid.UUID, code bool) {
	pInstance, _ := s.repo.Get(wId)
	pInstance.SetGeoStatus(code)
	pInstance.LiveData = false
	pInstance.GeoBrodacasting = false
	s.repo.Update(*pInstance)
}

func (s *PeripheralService) GetGeoLocation(wId uuid.UUID) domain.LastLocation {
	pInstance, _ := s.repo.Get(wId)
	return pInstance.GetGeoLocation()
}

func (s *PeripheralService) GetLiveSw(wId uuid.UUID) bool {
	pInstance, error := s.repo.Get(wId)
	if error != nil {
		return false
	}
	return pInstance.LiveData
}

func (s *PeripheralService) SetLiveSw(wId uuid.UUID, code bool) {
	pInstance, error := s.repo.Get(wId)
	if error != nil {
		log.Fatal("cannot set live sw")
		return
	}
	pInstance.LiveData = code
	s.repo.Update(*pInstance)
}
