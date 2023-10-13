package ports

import (
	"errors"

	"github.com/CAS735-F23/macrun-teamvs_/hrm/internal/core/domain"
)

var (
	ErrorListPlayersFailed  = errors.New("failed to list players")
	ErrorPlayerNotFound     = errors.New("the player session not found in repository")
	ErrorCreatePlayerFailed = errors.New("failed to add the player")
	ErrorUpdatePlayerFailed = errors.New("failed to update player")
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
