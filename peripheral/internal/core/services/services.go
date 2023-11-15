package services

import (
	"log"

	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/dto"
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/ports"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PeripheralService struct {
	repo ports.PeripheralRepository
}

func NewPeripheralService(repo ports.PeripheralRepository) *PeripheralService {
	return &PeripheralService{
		repo: repo,
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

func (s *PeripheralService) GetHRMAvgReading(hId uuid.UUID) (dto.LastHR, error) {
	pInstance, err := s.repo.GetByWorkoutId(hId)
	if err != nil {
		return dto.LastHR{}, err
	}
	return pInstance.GetAverageHRate(), nil
}

func (s *PeripheralService) GetHRMReading(hId uuid.UUID) (dto.LastHR, error) {
	pInstance, err := s.repo.GetByWorkoutId(hId)
	if err != nil {
		return dto.LastHR{}, err
	}
	return pInstance.GetHRate(), nil
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

func (s *PeripheralService) GetGeoLocation(wId uuid.UUID) (dto.LastLocation, error) {
	pInstance, err := s.repo.GetByWorkoutId(wId)
	if err != nil {
		return dto.LastLocation{}, err
	}
	return pInstance.GetGeoLocation(), nil
}

func (s *PeripheralService) GetLiveSw(wId uuid.UUID) (bool, error) {
	pInstance, err := s.repo.GetByWorkoutId(wId)
	if err != nil {
		return false, err
	}
	return pInstance.LiveStatus, nil
}

func (s *PeripheralService) SetLiveSw(wId uuid.UUID, code bool) error {
	pInstance, err := s.repo.GetByWorkoutId(wId)
	if err != nil {
		return err
	}
	pInstance.LiveStatus = code
	s.repo.Update(*pInstance)
	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
