package services

import (
	"log"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/ports"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PeripheralService struct {
	repo      ports.PeripheralRepository
	publisher ports.RabbitMQHandler
	client    ports.ZoneClient
}

func NewPeripheralService(repo ports.PeripheralRepository, handler ports.RabbitMQHandler, client ports.ZoneClient) *PeripheralService {
	return &PeripheralService{
		repo:      repo,
		publisher: handler,
		client:    client,
	}
}

func (s *PeripheralService) CreatePeripheral(pId uuid.UUID, hId uuid.UUID) error {
	p, err := domain.NewPeripheral(pId, hId, uuid.Nil, false, false)
	if err != nil {
		return err
	}
	s.repo.AddPeripheralIntance(p)
	return nil
}

func (s *PeripheralService) CheckStatusByHRMId(hId uuid.UUID) bool {
	pInstance, err := s.repo.GetByHRMId(hId)
	if err != nil {
		return false
	} else {
		s.repo.Update(*pInstance)
		return true
	}
}

func (s *PeripheralService) BindPeripheral(pId uuid.UUID, wId uuid.UUID, hId uuid.UUID, connected bool, sendToTrail bool) error {

	pInstance, err := s.repo.GetByHRMId(hId)
	if err != nil {

		p, _ := domain.NewPeripheral(pId, hId, wId, connected, sendToTrail)
		s.repo.AddPeripheralIntance(p)
		return nil
	}

	pInstance.PlayerId = pId
	pInstance.WorkoutId = wId
	pInstance.HRMDev.HRMStatus = connected
	pInstance.LiveStatus = sendToTrail
	s.repo.Update(*pInstance)
	return nil

}

func (s *PeripheralService) DisconnectPeripheral(wId uuid.UUID) error {
	err := s.repo.DeletePeripheralInstance(wId)
	return err
}

func (s *PeripheralService) GetHRMAvgReading(hId uuid.UUID) (uuid.UUID, time.Time, int, error) {
	pInstance, err := s.repo.GetByWorkoutId(hId)
	if err != nil {
		return uuid.Nil, time.Time{}, 0, err
	}
	return pInstance.HRMId, pInstance.HRMDev.HRateTime, pInstance.HRMDev.AverageHRate, nil
}

func (s *PeripheralService) GetHRMReading(hId uuid.UUID) (uuid.UUID, time.Time, int, error) {
	pInstance, err := s.repo.GetByWorkoutId(hId)
	if err != nil {
		return uuid.Nil, time.Time{}, 0, err
	}
	return pInstance.HRMId, pInstance.HRMDev.HRateTime, pInstance.HRMDev.HRate, nil
}

func (s *PeripheralService) SetHeartRateReading(hId uuid.UUID, reading int) error {
	pInstance, err := s.repo.GetByHRMId(hId)
	if err != nil {
		log.Fatal("error updating HRM", zap.Error(err))
		return err
	}
	pInstance.SetHRate(reading)
	s.repo.Update(*pInstance)
	return nil
}

func (s *PeripheralService) GetHRMDevStatus(wId uuid.UUID) (bool, error) {
	pInstance, err := s.repo.GetByWorkoutId(wId)
	if err != nil {
		return false, err
	}
	return pInstance.HRMDev.HRMStatus, nil
}

func (s *PeripheralService) SetHRMDevStatusByHRMId(hId uuid.UUID, code bool) error {
	pInstance, err := s.repo.GetByHRMId(hId)
	if err != nil {
		return err
	}
	pInstance.HRMDev.HRMStatus = code
	s.repo.Update(*pInstance)
	return nil
}

func (s *PeripheralService) SetHRMDevStatus(wId uuid.UUID, code bool) error {
	pInstance, err := s.repo.GetByWorkoutId(wId)
	if err != nil {
		return err
	}
	pInstance.HRMDev.HRMStatus = code
	s.repo.Update(*pInstance)
	return nil
}

func (s *PeripheralService) SetGeoLocation(wId uuid.UUID, longitude float64, latitude float64) error {
	pInstance, err := s.repo.GetByWorkoutId(wId)
	if err != nil {
		return err
	}
	pInstance.SetLocation(longitude, latitude)
	s.repo.Update(*pInstance)
	return nil
}

func (s *PeripheralService) GetGeoDevStatus(wId uuid.UUID) (bool, error) {
	pInstance, err := s.repo.GetByWorkoutId(wId)
	if err != nil {
		return false, err
	}
	return pInstance.GeoDev.GeoStatus, nil
}

func (s *PeripheralService) SetGeoDevStatus(wId uuid.UUID, code bool) error {
	pInstance, err := s.repo.GetByWorkoutId(wId)
	if err != nil {
		return err
	}
	pInstance.GeoDev.GeoStatus = code
	s.repo.Update(*pInstance)
	return nil
}

func (s *PeripheralService) GetGeoLocation(wId uuid.UUID) (time.Time, float64, float64, uuid.UUID, error) {
	pInstance, err := s.repo.GetByWorkoutId(wId)
	if err != nil {
		return time.Time{}, 0.0, 0.0, uuid.Nil, err
	}
	return pInstance.GeoDev.LocationTime, pInstance.GeoDev.Longitude, pInstance.GeoDev.Latitude, pInstance.WorkoutId, nil
}

func (s *PeripheralService) GetLiveStatus(wId uuid.UUID) (bool, error) {
	pInstance, err := s.repo.GetByWorkoutId(wId)
	if err != nil {
		return false, err
	}
	return pInstance.LiveStatus, nil
}

func (s *PeripheralService) SetLiveStatus(wId uuid.UUID, code bool) error {
	pInstance, err := s.repo.GetByWorkoutId(wId)
	if err != nil {
		return err
	}
	pInstance.LiveStatus = code
	s.repo.Update(*pInstance)
	return nil
}

func (s *PeripheralService) SendLastLocation(wId uuid.UUID, latitude float64, longitude float64, time time.Time) error {
	err := s.publisher.SendLastLocation(wId, latitude, longitude, time)
	if err != nil {
		return err
	}
	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func (s *PeripheralService) GetTrailLocationInfo(trailId uuid.UUID) (float64, float64, float64, float64, error) {
	return s.client.GetTrailLocation(trailId)
}
