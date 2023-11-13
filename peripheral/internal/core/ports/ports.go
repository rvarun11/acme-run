package ports

import (
	"errors"

	"github.com/CAS735-F23/macrun-teamvsl/hrm/internal/core/domain"
	"github.com/google/uuid"
)

var (
	ErrorListPeripheralSFailed  = errors.New("failed to list Peripherals")
	ErrorPeripheralNotFound     = errors.New("the Peripheral session not found in repository")
	ErrorCreatePeripheralFailed = errors.New("failed to add the Peripheral")
	ErrorUpdatePeripheralFailed = errors.New("failed to update Peripheral")
)

type PeripheralService interface {
	ConnectPeripheral(pId uuid.UUID) error
	DisconnectPeripheral(pId uuid.UUID) error
	BindPeripheralToWorkout(pId uuid.UUID, workout uuid.UUID)
	Get(pId uuid.UUID) (*domain.Peripheral, error)
	Get(hrmId uuid.UUID) int
	SendPeripheralId()
	SendPeripheralLocation()
	ReadPeripheralLocation(Longitude float64, Latitude float64)
	ReadPeripheralHeartRate(pId uuid.UUID, rate int)
}

type PeripheralRepository interface {
	AddPeripheralIntance(p domain.Peripheral) error
	DeletePeripheralInstance(pId uuid.UUID) error
	Get(pID uuid.UUID) (*domain.Peripheral, error)
	Update(p domain.Peripheral) error
	List() ([]*domain.Peripheral, error)
}
