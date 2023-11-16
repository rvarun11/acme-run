package ports

import (
	"errors"

	"github.com/CAS735-F23/macrun-teamvsl/trail/internal/core/domain"

	"github.com/google/uuid"
)

var (
	ErrorListWorkoutsFailed      = errors.New("failed to list workout")
	ErrorWorkoutNotFound         = errors.New("workout not found in repository")
	ErrorCreateWorkoutFailed     = errors.New("failed to create the workout")
	ErrorUpdateWorkoutFailed     = errors.New("failed to update workout")
	ErrorTrailManagerInvalidArgs = errors.New("invalid arguments in trail manager")
	ErrorInvalidTrail            = errors.New("invalid trail in creation")
)

type TrailManagerService interface {
	// retrieveLocation()
	// getCloestShelter()
	// GetShelterDistance() (float64, error)
	// GetTrail(id uuid.UUID) (*domain.Trail, error)
	// GetShelter(id uuid.UUID) (*domain.Shelter, error)
	// GetClosestShelter(currentLongitude, currentLatitude float64) (uuid.UUID, error)
	// GetClosestTrail(currentLongitude float64, currentLatitude float64) (uuid.UUID, error)
	// SelectTrail(wId uuid.UUID, tId uuid.UUID, option string) error
	// NewTrailManagerService(rTM MemoryRepository, rT TrailRepository, rS ShelterRepository) (*TrailManagerService, error)
	CreateTrailManager(wId uuid.UUID) error
	DisconnectTrailManager(wId uuid.UUID) error
	GetShelter(id uuid.UUID) (*domain.Shelter, error)
	GetTrail(id uuid.UUID) (*domain.Trail, error)
	calculateDistance(Longitude1 float64, Latitude1, Longitude2 float64, Latitude2 float64) (float64, error)
	SelectTrail(wId uuid.UUID, tId uuid.UUID, option string) error
	GetShelterDistance(wId uuid.UUID, tId uuid.UUID, sId uuid.UUID) (float64, error)
	GetTrailDistance(wId uuid.UUID, tId uuid.UUID, sId uuid.UUID) (float64, error)
	GetClosestShelter(currentLongitude, currentLatitude float64) (uuid.UUID, error)
	GetClosestTrail(currentLongitude float64, currentLatitude float64) (uuid.UUID, error)
}

type TrailRepository interface {
	CreateTrail(name string, startLat, startLong, endLat, endLong float64, closestShelterID uuid.UUID) (uuid.UUID, error)
	UpdateTrailByID(id uuid.UUID, name string, startLat, startLong, endLat, endLong float64, closestShelterID uuid.UUID) error
	DeleteTrailByID(id uuid.UUID) error
	GetTrailByID(id uuid.UUID) (domain.Trail, error)
}

type ShelterRepository interface {
	CreateShelter(name string, availability bool, lat, long float64) (uuid.UUID, error)
	UpdateShelterByID(id uuid.UUID, name string, availability bool, lat, long float64) error
	DeleteShelterByID(id uuid.UUID) error
	GetShelterByID(id uuid.UUID) (*domain.Shelter, error)
}

type TrailManagerRepository interface {
	GetByWorkoutId(wId uuid.UUID) (*domain.TrailManager, error)
	List() ([]*domain.TrailManager, error)
	Update(t domain.TrailManager) error
	DeleteTrailManagerInstance(wId uuid.UUID) error
	AddTrailManagerIntance(t domain.TrailManager) error
}
