package ports

import (
	"errors"

	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/dto"
	"github.com/google/uuid"
)

var (
	ErrorListPeripheralSFailed  = errors.New("failed to list Peripherals")
	ErrorPeripheralNotFound     = errors.New("the Peripheral session not found in repository")
	ErrorCreatePeripheralFailed = errors.New("failed to add the Peripheral")
	ErrorUpdatePeripheralFailed = errors.New("failed to update Peripheral")
	ErrorListPeripheralFailed   = errors.New("failed to list Peripheral")
)

// LS-TODO: Remove or comment out the unused functions

// LS-TODO: The sergvices should return the domain object with error
type PeripheralService interface {
	CreatePeripheral(pId uuid.UUID, hId uuid.UUID) error
	CheckStatusByHRMId(hId uuid.UUID) bool
	BindPeripheral(pId uuid.UUID, wId uuid.UUID, hId uuid.UUID, connected bool, sendToTrail bool) error
	DisconnectPeripheral(wId uuid.UUID) error
	GetHRMAvgReading(hId uuid.UUID) (dto.LastHR, error)
	GetHRMReading(hId uuid.UUID) (dto.LastHR, error)
	SetHeartRateReading(hId uuid.UUID, reading int) error
	GetHRMDevStatus(wId uuid.UUID) (bool, error)
	SetHRMDevStatusByHRMId(hId uuid.UUID, code bool) error
	SetHRMDevStatus(wId uuid.UUID, code bool) error
	SetGeoLocation(wId uuid.UUID, longitude float64, latitude float64) error
	GetGeoDevStatus(wId uuid.UUID) (bool, error)
	SetGeoDevStatus(wId uuid.UUID, code bool) error
	GetGeoLocation(wId uuid.UUID) (dto.LastLocation, error)
	GetLiveSw(wId uuid.UUID) (bool, error)
	SetLiveSw(wId uuid.UUID, code bool) error
}

type PeripheralRepository interface {
	AddPeripheralIntance(p domain.Peripheral) error
	DeletePeripheralInstance(pId uuid.UUID) error
	DeletePeripheralInstanceByHRMId(hId uuid.UUID) error
	GetByWorkoutId(wID uuid.UUID) (*domain.Peripheral, error)
	GetByPlayerId(pID uuid.UUID) (*domain.Peripheral, error)
	GetByHRMId(hID uuid.UUID) (*domain.Peripheral, error)
	Update(p domain.Peripheral) error
	List() ([]*domain.Peripheral, error)
}
