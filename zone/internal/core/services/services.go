package services

import (
	"fmt"
	"math"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/core/ports"
	logger "github.com/CAS735-F23/macrun-teamvsl/zone/log"
	"github.com/google/uuid"
	"github.com/umahmood/haversine"
	"go.uber.org/zap"
)

type ZoneService struct {
	repo                     ports.ZoneManagerRepository
	shelterDistancePublisher ports.ShelterDistancePublisher
}

func NewZoneService(repo ports.ZoneManagerRepository, shelterDistancePublisher ports.ShelterDistancePublisher) (*ZoneService, error) {
	return &ZoneService{
		repo:                     repo,
		shelterDistancePublisher: shelterDistancePublisher,
	}, nil
}

func (zs *ZoneService) CreateTrail(name string, zId uuid.UUID, startLatitude float64, startLongitude float64, endLatitude float64, endLongitude float64) (uuid.UUID, error) {
	res, err := zs.repo.CreateTrail(name, zId, startLatitude, startLongitude, endLatitude, endLongitude)
	if err != nil {
		return uuid.Nil, err
	}
	logger.Info("trail created successfully", zap.Any("trail_id", res))
	return res, nil
}

func (zs *ZoneService) UpdateTrail(tid uuid.UUID, name string, zId uuid.UUID, startLatitude float64, startLongitude float64, endLatitude float64, endLongitude float64) error {
	err := zs.repo.UpdateTrailByID(tid, name, zId, startLatitude, startLongitude, endLatitude, endLongitude)
	if err != nil {
		return err
	}
	return nil
}

func (zs *ZoneService) DeleteTrail(tId uuid.UUID) error {
	err := zs.repo.DeleteTrailByID(tId)
	if err != nil {
		return err
	}
	return nil
}

func (zs *ZoneService) GetTrailByID(id uuid.UUID) (*domain.Trail, error) {
	trail, err := zs.repo.GetTrailByID(id)
	return trail, err
}

func (zs *ZoneService) CheckTrail(id uuid.UUID) error {
	trail, err := zs.repo.GetTrailByID(id)
	if err != nil || trail.TrailID != id {
		return err
	}
	return nil
}

// VR-TODO This may be needed later
// func (zs *ZoneService) GetCurrentLocation(wId uuid.UUID) (float64, float64, error) {
// 	tmInstance, err := zs.repo.GetByWorkoutId(wId)
// 	if err != nil {
// 		return math.MaxFloat64, math.MaxFloat64, err
// 	}
// 	return tmInstance.CurrentLongitude, tmInstance.CurrentLatitude, nil
// }

func (zs *ZoneService) GetClosestTrail(zId uuid.UUID, currentLongitude float64, currentLatitude float64) (uuid.UUID, error) {

	trails, err := zs.repo.ListTrailsByZoneId(zId)
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
	logger.Info("closest trail info retrieved", zap.Any("trail", closestTrail.TrailID))
	return uuid.Nil, nil // Or return an appropriate error if necessary
}

func (zs *ZoneService) CreateShelter(name string, tId uuid.UUID, availability bool, lat, long float64) (uuid.UUID, error) {
	sId, err := zs.repo.CreateShelter(name, tId, availability, lat, long)
	if err != nil {
		return uuid.Nil, err
	} else {
		logger.Info("shelter created successfully", zap.Any("shelter_id", sId))
		return sId, nil
	}
}

func (zs *ZoneService) UpdateShelter(id uuid.UUID, name string, tId uuid.UUID, availability bool, lat, long float64) error {

	err := zs.repo.UpdateShelterByID(id, tId, name, availability, lat, long)
	if err != nil {
		logger.Error("Zone: failed to updater shelter", zap.Error(err))
		return err
	}
	return err

}

func (zs *ZoneService) DeleteShelter(id uuid.UUID) error {
	return zs.repo.DeleteShelterByID(id)
}

func (zs *ZoneService) GetShelterByID(id uuid.UUID) (*domain.Shelter, error) {
	shelter, err := zs.repo.GetShelterByID(id)
	if err != nil {
		logger.Error("Zone: shleter location can't be retrieved", zap.Error(err))
	}
	return shelter, err
}

func (zs *ZoneService) CheckShelter(id uuid.UUID) error {
	s, err := zs.repo.GetShelterByID(id)
	if err != nil || s.ShelterID != id {
		return err
	}
	return nil
}

func (zs *ZoneService) CreateZone(zName string) (uuid.UUID, error) {
	zId, err := zs.repo.CreateZone(zName)
	if err != nil {
		return uuid.Nil, err
	}
	logger.Info("zone created successfully", zap.Any("zone_id", zId))
	return zId, nil
}

func (zs *ZoneService) CheckZone(zId uuid.UUID) error {
	z, err := zs.repo.GetZoneByID(zId)

	if err != nil || z.ZoneID != zId {
		fmt.Println("z.ZoneID")
		return err
	}
	return nil
}

func (zs *ZoneService) CheckZoneByName(name string) error {
	z, err := zs.repo.GetZoneByName(name)
	if err != nil || z.ZoneName != name {
		return err
	}
	return nil
}

func (zs *ZoneService) UpdateZone(zId uuid.UUID, zName string) error {
	err := zs.repo.UpdateZone(zId, zName)
	if err != nil {
		return err
	}
	return nil
}

func (zs *ZoneService) DeleteZone(zId uuid.UUID) error {

	err := zs.CheckZone(zId)
	if err != nil {
		return nil
	}

	err = zs.repo.DeleteZone(zId)
	if err != nil {
		return err
	}
	return nil
}

func (zs *ZoneService) UpdateCurrentLocation(wId uuid.UUID, latitude float64, longitude float64, time time.Time) error {

	// Now push the shelter data data to the queue to the workout
	shelterId, distance, availability, _, err := zs.GetClosestShelter(longitude, latitude, time)
	if err != nil {
		logger.Error("error when getting cloest shelter info", zap.Error(err))
		return err
	}
	closestShelter, _ := zs.repo.GetShelterByID(shelterId)
	err = zs.shelterDistancePublisher.PublishShelterDistance(wId, shelterId, closestShelter.ShelterName, availability, distance)

	if err != nil {
		logger.Error("error when publishing shelter info", zap.Error(err))
	}
	logger.Debug("publishing shelter data to workout thru queue")

	return nil
}

func (zs *ZoneService) GetClosestShelter(longitude float64, latitude float64, time time.Time) (uuid.UUID, float64, bool, time.Time, error) {

	shelters, err := zs.repo.ListShelters()
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

func (zs *ZoneService) GetClosestShelterInfo(latitude float64, longitude float64) (uuid.UUID, float64, error) {
	shelters, err := zs.repo.ListShelters()
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
func (zs *ZoneService) ListZones() ([]*domain.Zone, error) {
	zones, err := zs.repo.ListZones()
	if err != nil {
		return nil, err
	}
	return zones, nil
}
