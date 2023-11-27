package services

import (
	"fmt"
	"math"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/core/ports"
	"github.com/CAS735-F23/macrun-teamvsl/zone/log"
	"github.com/google/uuid"
	"github.com/umahmood/haversine"
	"go.uber.org/zap"
)

type ZoneManagerService struct {
	repoTrail       ports.TrailRepository
	repoShelter     ports.ShelterRepository
	repoZone        ports.ZoneRepository
	repoZoneManager ports.ZoneManagerRepository
	publisher       ports.AMQPPublisher
}

// Factory for creating a new ZoneManager

func NewZoneManagerService(rTM ports.ZoneManagerRepository, rT ports.TrailRepository, rS ports.ShelterRepository, rZ ports.ZoneRepository, publisher ports.AMQPPublisher) (*ZoneManagerService, error) {
	return &ZoneManagerService{
		repoZoneManager: rTM,
		repoTrail:       rT,
		repoShelter:     rS,
		repoZone:        rZ,
		publisher:       publisher,
	}, nil
}

func (t *ZoneManagerService) CreateZoneManager(wId uuid.UUID) (uuid.UUID, error) {
	tm, err := domain.NewZoneManager(wId)
	if err != nil {
		return uuid.Nil, err
	}
	t.repoZoneManager.AddZoneManagerIntance(tm)
	return tm.ZoneManagerID, nil
}

func (t *ZoneManagerService) CloseZoneManager(wId uuid.UUID) error {
	return t.repoZoneManager.DeleteZoneManagerInstance(wId)
}

func (t *ZoneManagerService) GetZoneManagerByWorkoutId(id uuid.UUID) (*domain.ZoneManager, error) {
	tm, err := t.repoZoneManager.GetByWorkoutId(id)
	return tm, err
}

func (t *ZoneManagerService) CreateTrail(name string, zId uuid.UUID, startLatitude float64, startLongitude float64, endLatitude float64, endLongitude float64) (uuid.UUID, error) {
	res, err := t.repoTrail.CreateTrail(name, zId, startLatitude, startLongitude, endLatitude, endLongitude)
	if err != nil {
		return uuid.Nil, err
	}
	return res, nil
}

func (t *ZoneManagerService) UpdateTrail(tid uuid.UUID, name string, zId uuid.UUID, startLatitude float64, startLongitude float64, endLatitude float64, endLongitude float64) error {
	err := t.repoTrail.UpdateTrailByID(tid, name, zId, startLatitude, startLongitude, endLatitude, endLongitude)
	if err != nil {
		return err
	}
	return nil
}

func (t *ZoneManagerService) DeleteTrail(tId uuid.UUID) error {
	err := t.repoTrail.DeleteTrailByID(tId)
	log.Debug("deelte trail err", zap.Error(err))
	if err != nil {
		return err
	}
	return nil
}

func (t *ZoneManagerService) DisconnectZoneManager(wId uuid.UUID) error {
	err := t.repoZoneManager.DeleteZoneManagerInstance(wId)
	return err
}

func (t *ZoneManagerService) GetTrailByID(id uuid.UUID) (*domain.Trail, error) {
	trail, err := t.repoTrail.GetTrailByID(id)
	return trail, err
}

func (t *ZoneManagerService) CheckTrail(id uuid.UUID) error {
	trail, err := t.repoTrail.GetTrailByID(id)
	if err != nil || trail.TrailID != id {
		return err
	}
	return nil
}

func (t *ZoneManagerService) GetCurrentLocation(wId uuid.UUID) (float64, float64, error) {
	tmInstance, err := t.repoZoneManager.GetByWorkoutId(wId)
	if err != nil {
		return math.MaxFloat64, math.MaxFloat64, err
	}
	return tmInstance.CurrentLongitude, tmInstance.CurrentLatitude, nil
}

func (t *ZoneManagerService) GetClosestTrail(zId uuid.UUID, currentLongitude float64, currentLatitude float64) (uuid.UUID, error) {

	trails, err := t.repoTrail.ListTrailsByZoneId(zId)
	if err != nil {
		return uuid.Nil, err // Handle the error, possibly no trails available or DB error
	}
	var closestTrail *domain.Trail
	minDistance := math.MaxFloat64 // Initialize with the maximum float value
	distance := math.MaxFloat64
	for _, trail := range trails {
		point1 := haversine.Coord{Lat: currentLatitude, Lon: currentLongitude}
		point2 := haversine.Coord{Lat: trail.StartLatitude, Lon: trail.StartLongitude}
		_, distance = haversine.Distance(point1, point2)
		if distance < minDistance {
			minDistance = distance
			closestTrail = trail
		}
	}

	// If a closest trail is found, update the ZoneManager
	if closestTrail != nil {

		return closestTrail.TrailID, nil
	}

	return uuid.Nil, nil // Or return an appropriate error if necessary
}

func (t *ZoneManagerService) CreateShelter(name string, tId uuid.UUID, availability bool, lat, long float64) (uuid.UUID, error) {
	sId, err := t.repoShelter.CreateShelter(name, tId, availability, lat, long)
	if err != nil {
		return uuid.Nil, err
	} else {
		return sId, nil
	}
}

func (t *ZoneManagerService) UpdateShelter(id uuid.UUID, name string, tId uuid.UUID, availability bool, lat, long float64) error {
	return t.repoShelter.UpdateShelterByID(id, tId, name, availability, lat, long)
}

func (t *ZoneManagerService) DeleteShelter(id uuid.UUID) error {
	return t.repoShelter.DeleteShelterByID(id)
}

func (t *ZoneManagerService) GetShelterByID(id uuid.UUID) (*domain.Shelter, error) {
	shelter, err := t.repoShelter.GetShelterByID(id)
	return shelter, err
}

func (t *ZoneManagerService) CheckShelter(id uuid.UUID) error {
	s, err := t.repoShelter.GetShelterByID(id)
	if err != nil || s.ShelterID != id {
		return err
	}
	return nil
}

func (t *ZoneManagerService) CreateZone(zName string) (uuid.UUID, error) {
	zId, err := t.repoZone.CreateZone(zName)
	if err != nil {
		return uuid.Nil, err
	}
	return zId, nil
}

func (t *ZoneManagerService) CheckZone(zId uuid.UUID) error {
	z, err := t.repoZone.GetZoneByID(zId)

	if err != nil || z.ZoneID != zId {
		fmt.Println("z.ZoneID")
		return err
	}
	return nil
}

func (t *ZoneManagerService) CheckZoneByName(name string) error {
	z, err := t.repoZone.GetZoneByName(name)
	if err != nil || z.ZoneName != name {
		return err
	}
	return nil
}

func (t *ZoneManagerService) UpdateZone(zId uuid.UUID, zName string) error {
	err := t.repoZone.UpdateZone(zId, zName)
	if err != nil {
		return err
	}
	return nil
}

func (t *ZoneManagerService) DeleteZone(zId uuid.UUID) error {

	err := t.CheckZone(zId)
	if err != nil {
		return nil
	}

	err = t.repoZone.DeleteZone(zId)
	if err != nil {
		return err
	}
	return nil
}

func (t *ZoneManagerService) UpdateCurrentLocation(wId uuid.UUID, latitude float64, longitude float64, time time.Time) error {

	// Now push the shelter data data to the queue to the workout
	shelterId, distance, availability, _, err := t.GetClosestShelter(longitude, latitude, time)
	if err != nil {
		log.Error("error when getting cloest shelter info", zap.Error(err))
		return err
	}
	closestShelter, _ := t.GetShelterByID(shelterId)
	err = t.publisher.PublishShelterInfo(wId, shelterId, closestShelter.ShelterName, availability, distance)

	if err != nil {
		log.Error("error when publishing shelter info", zap.Error(err))
	}
	log.Debug("publishing shelter data to workout thru queue")

	return nil
}

func (t *ZoneManagerService) GetClosestShelter(longitude float64, latitude float64, time time.Time) (uuid.UUID, float64, bool, time.Time, error) {

	shelters, err := t.repoShelter.List()
	if err != nil {
		return uuid.Nil, math.MaxFloat64, false, time, err
	}

	var closestShelter *domain.Shelter
	minDistance := math.MaxFloat64 // Initialize with the maximum float value

	for _, shelter := range shelters {

		distance := 0.0
		point1 := haversine.Coord{Lat: latitude, Lon: longitude}
		point2 := haversine.Coord{Lat: shelter.Latitude, Lon: shelter.Longitude}
		_, distance = haversine.Distance(point1, point2)

		if distance < minDistance {
			minDistance = distance
			closestShelter = shelter
		}
	}

	// If a closest trail is found, update the ZoneManager
	if closestShelter != nil {

		return closestShelter.ShelterID, minDistance, closestShelter.ShelterAvailability, time, nil
	}

	return uuid.Nil, math.MaxFloat64, false, time, nil // Or return an appropriate error if necessary
}

func (t *ZoneManagerService) GetClosestShelterInfo(latitude float64, longitude float64) (uuid.UUID, float64, error) {
	shelters, err := t.repoShelter.List()
	if err != nil || len(shelters) == 0 {
		return uuid.Nil, math.MaxFloat64, err
	}
	var closestShelter *domain.Shelter
	minDistance := math.MaxFloat64 // Initialize with the maximum float value

	for _, shelter := range shelters {

		distance := 0.0
		point1 := haversine.Coord{Lat: latitude, Lon: longitude}
		point2 := haversine.Coord{Lat: shelter.Latitude, Lon: shelter.Longitude}
		_, distance = haversine.Distance(point1, point2)

		if distance < minDistance {
			minDistance = distance
			closestShelter = shelter
		}
	}
	return closestShelter.ShelterID, minDistance, nil

}

// Add this function to the ZoneManagerService type
func (t *ZoneManagerService) ListZones() ([]*domain.Zone, error) {
	zones, err := t.repoZone.List()
	if err != nil {
		return nil, err
	}
	return zones, nil
}
