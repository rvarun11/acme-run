package ports

import (
	"errors"

	"github.com/CAS735-F23/macrun-teamvsl/trail/internal/core/domain"
	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

var (
	ErrorCreateTrailManagerFailed = errors.New("failed to create the trail manager in repo")
	ErrorUpdateTrailManagerFailed = errors.New("failed to update trail manager in repo")
	ErrorTrailManagerInvalidArgs  = errors.New("invalid arguments in trail manager in repo")
	ErrorInvalidTrail             = errors.New("invalid trail")
	ErrorInvalidShelter           = errors.New("invalid shelter")
	ErrorTrailManagerlNotFound    = errors.New("trail manager not found")
	ErrorListTrailManagerFailed   = errors.New("listing trails manager failed")
)

type TrailManagerService interface {
	CreateTrailManager(wId uuid.UUID) error
	GetShelterByID(id uuid.UUID) (*domain.Shelter, error)
	GetTrail(id uuid.UUID) (*domain.Trail, error)
	CalculateDistance(Longitude1 float64, Latitude1, Longitude2 float64, Latitude2 float64) (float64, error)
	GetShelterDistance(wId uuid.UUID, tId uuid.UUID, sId uuid.UUID) (float64, error)
	GetTrailDistance(wId uuid.UUID, tId uuid.UUID, sId uuid.UUID) (float64, error)
	GetClosestShelter(currentLongitude, currentLatitude float64) (uuid.UUID, error)
	GetClosestTrail(zId uuid.UUID, currentLongitude float64, currentLatitude float64) (uuid.UUID, error)
	SetCurrentLocation(wId uuid.UUID, longitude float64, latitude float64)
	CreateTrail(tid uuid.UUID, name string, startLatitude float64, startLongitude float64, endLatitude float64, endLongitude float64, shelterId uuid.UUID) (uuid.UUID, error)
	GetTrailInfo(ctx *gin.Context)
}

type TrailRepository interface {
	CreateTrail(name string, zId uuid.UUID, startLat, startLong, endLat, endLong float64) (uuid.UUID, error)
	UpdateTrailByID(id uuid.UUID, name string, zId uuid.UUID, startLat, startLong, endLat, endLong float64) error
	DeleteTrailByID(id uuid.UUID) error
	GetTrailByID(id uuid.UUID) (*domain.Trail, error)
	List() ([]*domain.Trail, error)
	ListTrailsByZoneId(zId uuid.UUID) ([]*domain.Trail, error)
}

type ShelterRepository interface {
	CreateShelter(name string, tId uuid.UUID, availability bool, lat, long float64) (uuid.UUID, error)
	UpdateShelterByID(id uuid.UUID, tId uuid.UUID, name string, availability bool, lat, long float64) error
	DeleteShelterByID(id uuid.UUID) error
	GetShelterByID(id uuid.UUID) (*domain.Shelter, error)
	List() ([]*domain.Shelter, error)
	ListSheltersByTrailId(tId uuid.UUID) ([]*domain.Shelter, error)
}

type ZoneRepository interface {
	CreateZone(name string) (uuid.UUID, error)
	UpdateZone(id uuid.UUID, name string) error
	DeleteZone(id uuid.UUID) error
	GetZoneByID(id uuid.UUID) (*domain.Zone, error)
	GetZoneByName(name string) (*domain.Zone, error)
	List() ([]*domain.Zone, error)
}

type TrailManagerRepository interface {
	GetByWorkoutId(wId uuid.UUID) (*domain.TrailManager, error)
	Update(t domain.TrailManager) error
	DeleteTrailManagerInstance(wId uuid.UUID) error
	AddTrailManagerIntance(t domain.TrailManager) error
}

type AMQPPublisher interface {
	PublishShelterInfo(sId uuid.UUID, name string, availability bool, distance float64) error
}
