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

type ZoneService struct {
	repoZoneManager ports.ZoneManagerRepository
	repoDB          ports.DBRepository
	publisher       ports.AMQPPublisher
}

func NewZoneService(rTM ports.ZoneManagerRepository, rDB ports.DBRepository, publisher ports.AMQPPublisher) (*ZoneService, error) {
	return &ZoneService{
		repoZoneManager: rTM,
		repoDB:          rDB,
		publisher:       publisher,
	}, nil
}

func (t *ZoneService) CreateZoneManager(wId uuid.UUID) (uuid.UUID, error) {
	tm, err := domain.NewZoneManager(wId)
	if err != nil {
		return uuid.Nil, err
	}
	t.repoZoneManager.AddZoneManagerIntance(tm)
	return tm.ZoneManagerID, nil
}

func (t *ZoneService) CloseZoneManager(wId uuid.UUID) error {
	return t.repoZoneManager.DeleteZoneManagerInstance(wId)
}

func (t *ZoneService) GetZoneManagerByWorkoutId(id uuid.UUID) (*domain.ZoneManager, error) {
	tm, err := t.repoZoneManager.GetByWorkoutId(id)
	return tm, err
}

func (t *ZoneService) CreateTrail(name string, zId uuid.UUID, startLatitude float64, startLongitude float64, endLatitude float64, endLongitude float64) (uuid.UUID, error) {
	res, err := t.repoDB.CreateTrail(name, zId, startLatitude, startLongitude, endLatitude, endLongitude)
	if err != nil {
		return uuid.Nil, err
	}
	log.Info("Zone: trail created")
	return res, nil
}

func (t *ZoneService) UpdateTrail(tid uuid.UUID, name string, zId uuid.UUID, startLatitude float64, startLongitude float64, endLatitude float64, endLongitude float64) error {
	err := t.repoDB.UpdateTrailByID(tid, name, zId, startLatitude, startLongitude, endLatitude, endLongitude)
	if err != nil {
		return err
	}
	log.Info("Zone: trail updated")
	return nil
}

func (t *ZoneService) DeleteTrail(tId uuid.UUID) error {
	err := t.repoDB.DeleteTrailByID(tId)
	if err != nil {
		return err
	}
	log.Info("Zone: trail deleted")
	return nil
}

func (t *ZoneService) DisconnectZoneManager(wId uuid.UUID) error {
	err := t.repoZoneManager.DeleteZoneManagerInstance(wId)
	return err
}

func (t *ZoneService) GetTrailByID(id uuid.UUID) (*domain.Trail, error) {
	trail, err := t.repoDB.GetTrailByID(id)
	log.Info("Zone: trail detailed info retrieved")
	return trail, err
}

func (t *ZoneService) CheckTrail(id uuid.UUID) error {
	trail, err := t.repoDB.GetTrailByID(id)
	if err != nil || trail.TrailID != id {
		return err
	}
	return nil
}

func (t *ZoneService) GetCurrentLocation(wId uuid.UUID) (float64, float64, error) {
	tmInstance, err := t.repoZoneManager.GetByWorkoutId(wId)
	if err != nil {
		return math.MaxFloat64, math.MaxFloat64, err
	}
	return tmInstance.CurrentLongitude, tmInstance.CurrentLatitude, nil
}

func (t *ZoneService) GetClosestTrail(zId uuid.UUID, currentLongitude float64, currentLatitude float64) (uuid.UUID, error) {

	trails, err := t.repoDB.ListTrailsByZoneId(zId)
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
		log.Info("Zone: cloest shelter id retrieved")
		return closestTrail.TrailID, nil
	}

	return uuid.Nil, nil // Or return an appropriate error if necessary
}

func (t *ZoneService) CreateShelter(name string, tId uuid.UUID, availability bool, lat, long float64) (uuid.UUID, error) {
	sId, err := t.repoDB.CreateShelter(name, tId, availability, lat, long)
	if err != nil {
		return uuid.Nil, err
	} else {
		log.Info("Zone: shelter created")
		return sId, nil
	}
}

func (t *ZoneService) UpdateShelter(id uuid.UUID, name string, tId uuid.UUID, availability bool, lat, long float64) error {

	err := t.repoDB.UpdateShelterByID(id, tId, name, availability, lat, long)
	if err != nil {
		log.Error("Zone: failed to updater shelter", zap.Error(err))
		return err
	}
	log.Info("Zone: shelter updated")
	return err

}

func (t *ZoneService) DeleteShelter(id uuid.UUID) error {
	return t.repoDB.DeleteShelterByID(id)
}

func (t *ZoneService) GetShelterByID(id uuid.UUID) (*domain.Shelter, error) {
	shelter, err := t.repoDB.GetShelterByID(id)
	if err == nil {
		log.Info("Zone: shelter location info retrieved")
	} else {
		log.Error("Zone: shleter location can't be retrieved", zap.Error(err))
	}
	return shelter, err
}

func (t *ZoneService) CheckShelter(id uuid.UUID) error {
	s, err := t.repoDB.GetShelterByID(id)
	if err != nil || s.ShelterID != id {
		return err
	}
	return nil
}

func (t *ZoneService) CreateZone(zName string) (uuid.UUID, error) {
	zId, err := t.repoDB.CreateZone(zName)
	if err != nil {
		return uuid.Nil, err
	}
	log.Info("Zone: zone created")
	return zId, nil
}

func (t *ZoneService) CheckZone(zId uuid.UUID) error {
	z, err := t.repoDB.GetZoneByID(zId)

	if err != nil || z.ZoneID != zId {
		fmt.Println("z.ZoneID")
		return err
	}
	return nil
}

func (t *ZoneService) CheckZoneByName(name string) error {
	z, err := t.repoDB.GetZoneByName(name)
	if err != nil || z.ZoneName != name {
		return err
	}
	return nil
}

func (t *ZoneService) UpdateZone(zId uuid.UUID, zName string) error {
	err := t.repoDB.UpdateZone(zId, zName)
	if err != nil {
		return err
	}
	log.Info("Zone: shelter updated")
	return nil
}

func (t *ZoneService) DeleteZone(zId uuid.UUID) error {

	err := t.CheckZone(zId)
	if err != nil {
		return nil
	}

	err = t.repoDB.DeleteZone(zId)
	if err != nil {
		return err
	}
	log.Info("Zone: zone deleted")
	return nil
}

func (t *ZoneService) UpdateCurrentLocation(wId uuid.UUID, latitude float64, longitude float64, time time.Time) error {

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

func (t *ZoneService) GetClosestShelter(longitude float64, latitude float64, time time.Time) (uuid.UUID, float64, bool, time.Time, error) {

	shelters, err := t.repoDB.ListShelters()
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
		log.Info("Zone: cloest shelter info retrieved")
		return closestShelter.ShelterID, minDistance, closestShelter.ShelterAvailability, time, nil
	}

	return uuid.Nil, math.MaxFloat64, false, time, nil // Or return an appropriate error if necessary
}

func (t *ZoneService) GetClosestShelterInfo(latitude float64, longitude float64) (uuid.UUID, float64, error) {
	shelters, err := t.repoDB.ListShelters()
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

// Add this function to the ZoneService type
func (t *ZoneService) ListZones() ([]*domain.Zone, error) {
	zones, err := t.repoDB.ListZones()
	if err != nil {
		return nil, err
	}
	return zones, nil
}
