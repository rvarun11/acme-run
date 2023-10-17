package ports

import (
	"errors"

	"github.com/CAS735-F23/macrun-teamvs_/hrm/internal/core/domain"
)

var (
	ErrorListHRMSFailed  = errors.New("failed to list hrms")
	ErrorHRMNotFound     = errors.New("the hrm session not found in repository")
	ErrorCreateHRMFailed = errors.New("failed to add the hrm")
	ErrorUpdateHRMFailed = errors.New("failed to update hrm")
)

type HRMService interface {
	List() ([]*domain.HRM, error)
	Get(id string) (*domain.HRM, error)
	Create(hrm domain.HRM) error
	Update(hrm domain.HRM) (*domain.HRM, error)
}

type HRMRepository interface {
	List() ([]*domain.HRM, error)
	Create(hrm domain.HRM) error
	Get(id string) (*domain.HRM, error)
	Update(hrm domain.HRM) error
}
