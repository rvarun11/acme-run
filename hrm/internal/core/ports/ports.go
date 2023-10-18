package ports

import (
	"errors"

	"github.com/CAS735-F23/macrun-teamvsl/hrm/internal/core/domain"
	"github.com/google/uuid"
)

var (
	ErrorListHRMSFailed  = errors.New("failed to list hrms")
	ErrorHRMNotFound     = errors.New("the hrm session not found in repository")
	ErrorCreateHRMFailed = errors.New("failed to add the hrm")
	ErrorUpdateHRMFailed = errors.New("failed to update hrm")
)

type HRMService interface {
	ConnectHRM(HRMId uuid.UUID) error
	DisconnectHRM(HRMId uuid.UUID) error
	BindHRMtoWorkout(HRMId uuid.UUID, workout uuid.UUID)
	Get(hrmId uuid.UUID) (*domain.HRM, error)
	SendHRM()
}

type HRMRepository interface {
	AddHRMIntance(hrm domain.HRM) error
	DeleteHRMInstance(hrmId uuid.UUID) error
	Get(hrmID uuid.UUID) (*domain.HRM, error)
	Update(hrm domain.HRM) error
	List() ([]*domain.HRM, error)
}
