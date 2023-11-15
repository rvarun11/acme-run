package ports

import (
	"errors"

	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/domain"
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
	ConnectPeripheral(pId uuid.UUID) error
	DisconnectPeripheral(pId uuid.UUID) error
	BindPeripheralToWorkout(pId uuid.UUID, workout uuid.UUID)
	Get(pId uuid.UUID) (*domain.Peripheral, error)
	SetHRMReading(wId uuid.UUID, reading int)
	// SendPeripheralId()
	GetHRMReading(wId uuid.UUID) *domain.LastHR
	GetGeoLocation(wId uuid.UUID) *domain.LastLocation
	SendPeripheralLocation()

	ReadPeripheralLocation(Longitude float64, Latitude float64)
	ReadPeripheralHeartRate(pId uuid.UUID, rate int)
}

type PeripheralRepository interface {
	AddPeripheralIntance(p domain.Peripheral) error
	DeletePeripheralInstance(pId uuid.UUID) error
	DeletePeripheralInstanceByHRMId(hId uuid.UUID) error
	Get(wID uuid.UUID) (*domain.Peripheral, error)
	GetByPlayerId(pID uuid.UUID) (*domain.Peripheral, error)
	GetByHRMId(hID uuid.UUID) (*domain.Peripheral, error)
	Update(p domain.Peripheral) error
	List() ([]*domain.Peripheral, error)
}
