package services

import (
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/ports"
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/log"
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
	p, err := domain.NewPeripheral(pId, hId, uuid.Nil, false, false, false)
	if err != nil {
		return ports.ErrorCreatePeripheralFailed
	}
	s.repo.AddPeripheralIntance(p)
	log.Debug("peripheral instance created")
	return nil
}

func (s *PeripheralService) CheckStatusByHRMId(hId uuid.UUID) bool {
	pInstance, err := s.repo.GetByHRMId(hId)
	if err != nil {
		return false
	} else {
		s.repo.Update(pInstance)
		return true
	}
}

func (s *PeripheralService) BindPeripheral(pId uuid.UUID, wId uuid.UUID, hId uuid.UUID, connected bool, toShelter bool) error {

	pInstance, err := s.repo.GetByHRMId(hId)
	if err != nil {
		log.Debug("creating a instance")
		p, _ := domain.NewPeripheral(pId, hId, wId, connected, true, toShelter)
		s.repo.AddPeripheralIntance(p)
		pInstance, _ = s.repo.GetByHRMId(hId)
	}

	pInstance.PlayerId = pId
	pInstance.WorkoutId = wId
	pInstance.HRMDev.HRMStatus = connected
	pInstance.ToShelter = toShelter
	s.repo.Update(pInstance)
	log.Debug("peripheral binding success", zap.Any("created", true))
	return nil

}

func (s *PeripheralService) DisconnectPeripheral(wId uuid.UUID) error {
	err := s.repo.DeletePeripheralInstance(wId)
	if err != nil {
		return ports.ErrorPeripheralNotFound
	}
	return err
}

func (s *PeripheralService) GetHRMAvgReading(hId uuid.UUID) (uuid.UUID, time.Time, int, error) {
	pInstance, err := s.repo.GetByWorkoutId(hId)
	if err != nil || !pInstance.HRMDev.HRMStatus {
		return uuid.Nil, time.Time{}, 0, ports.ErrorPeripheralNotFound
	}

	return pInstance.HRMId, pInstance.HRMDev.HRateTime, pInstance.HRMDev.AverageHRate, nil
}

func (s *PeripheralService) GetHRMReading(wId uuid.UUID) (uuid.UUID, time.Time, int, error) {
	pInstance, err := s.repo.GetByWorkoutId(wId)
	if err != nil {
		return uuid.Nil, time.Time{}, 0, ports.ErrorPeripheralNotFound
	}
	return pInstance.HRMId, pInstance.HRMDev.HRateTime, pInstance.HRMDev.HRate, nil
}

func (s *PeripheralService) SetHeartRateReading(hId uuid.UUID, reading int) error {
	pInstance, err := s.repo.GetByHRMId(hId)
	if err != nil {
		log.Fatal("error updating HRM", zap.Error(err))
		return ports.ErrorPeripheralNotFound
	}
	pInstance.SetHRate(reading)
	s.repo.Update(pInstance)
	return nil
}

func (s *PeripheralService) GetHRMDevStatus(wId uuid.UUID) (bool, error) {
	pInstance, err := s.repo.GetByWorkoutId(wId)
	if err != nil {
		return false, ports.ErrorPeripheralNotFound
	}
	return pInstance.HRMDev.HRMStatus, nil
}

func (s *PeripheralService) GetHRMDevStatusByPlayerId(pId uuid.UUID) (bool, error) {
	pInstance, err := s.repo.GetByPlayerId(pId)
	if err != nil {
		return false, ports.ErrorPeripheralNotFound
	}
	return pInstance.HRMDev.HRMStatus, nil
}

func (s *PeripheralService) GetHRMDevStatusByHRMId(hId uuid.UUID) (bool, error) {
	pInstance, err := s.repo.GetByPlayerId(hId)
	if err != nil {
		return false, ports.ErrorPeripheralNotFound
	}
	return pInstance.HRMDev.HRMStatus, nil
}

func (s *PeripheralService) SetHRMDevStatusByHRMId(hId uuid.UUID, code bool) error {
	pInstance, err := s.repo.GetByHRMId(hId)
	if err != nil {
		return ports.ErrorPeripheralNotFound
	}
	pInstance.HRMDev.HRMStatus = code
	s.repo.Update(pInstance)
	return nil
}

func (s *PeripheralService) SetHRMDevStatus(wId uuid.UUID, code bool) error {
	pInstance, err := s.repo.GetByWorkoutId(wId)
	if err != nil {
		return ports.ErrorPeripheralNotFound
	}
	pInstance.HRMDev.HRMStatus = code
	s.repo.Update(pInstance)
	return nil
}

func (s *PeripheralService) SetGeoLocation(wId uuid.UUID, longitude float64, latitude float64) error {
	pInstance, err := s.repo.GetByWorkoutId(wId)
	if err != nil {
		return ports.ErrorPeripheralNotFound
	}
	pInstance.SetLocation(longitude, latitude)
	s.repo.Update(pInstance)
	return nil
}

func (s *PeripheralService) GetGeoDevStatus(wId uuid.UUID) (bool, error) {
	pInstance, err := s.repo.GetByWorkoutId(wId)
	if err != nil {
		return false, ports.ErrorPeripheralNotFound
	}
	return pInstance.GeoDev.GeoStatus, nil
}

func (s *PeripheralService) SetGeoDevStatus(wId uuid.UUID, code bool) error {
	pInstance, err := s.repo.GetByWorkoutId(wId)
	if err != nil {
		return ports.ErrorPeripheralNotFound
	}
	pInstance.GeoDev.GeoStatus = code
	s.repo.Update(pInstance)
	return nil
}

func (s *PeripheralService) GetGeoLocation(wId uuid.UUID) (time.Time, float64, float64, uuid.UUID, error) {
	pInstance, err := s.repo.GetByWorkoutId(wId)
	if err != nil {
		return time.Time{}, 0.0, 0.0, uuid.Nil, ports.ErrorPeripheralNotFound
	}
	return pInstance.GeoDev.LocationTime, pInstance.GeoDev.Longitude, pInstance.GeoDev.Latitude, pInstance.WorkoutId, nil
}

func (s *PeripheralService) GetLiveStatus(wId uuid.UUID) (bool, error) {
	pInstance, err := s.repo.GetByWorkoutId(wId)
	if err != nil {
		return false, ports.ErrorPeripheralNotFound
	}
	return pInstance.LiveStatus, nil
}

func (s *PeripheralService) SetLiveStatus(wId uuid.UUID, code bool) error {
	pInstance, err := s.repo.GetByWorkoutId(wId)
	if err != nil {
		return ports.ErrorPeripheralNotFound
	}
	pInstance.LiveStatus = code
	s.repo.Update(pInstance)
	return nil
}

func (s *PeripheralService) SendLastLocation(wId uuid.UUID, latitude float64, longitude float64, time time.Time) error {
	pInstance, _ := s.repo.GetByWorkoutId(wId)
	err := s.publisher.SendLastLocation(wId, latitude, longitude, time, pInstance.ToShelter)
	if err != nil {
		return ports.ErrorPeripheralPublishFailed
	}
	return nil
}

func (s *PeripheralService) GetTrailLocationInfo(trailId uuid.UUID) (float64, float64, float64, float64, error) {
	return s.client.GetTrailLocation(trailId)
}
