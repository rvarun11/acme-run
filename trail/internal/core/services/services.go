package services

import (
	"math"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/trail/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/trail/internal/core/ports"
	"github.com/CAS735-F23/macrun-teamvsl/trail/log"
	"github.com/google/uuid"
	"github.com/umahmood/haversine"
	"go.uber.org/zap"
)

type TrailManagerService struct {
	repoT     ports.TrailRepository
	repoS     ports.ShelterRepository
	repoZ     ports.ZoneRepository
	repoTM    ports.TrailManagerRepository
	publisher ports.AMQPPublisher
}

// Factory for creating a new TrailManager

func NewTrailManagerService(rTM ports.TrailManagerRepository, rT ports.TrailRepository, rS ports.ShelterRepository, rZ ports.ZoneRepository, publisher ports.AMQPPublisher) (*TrailManagerService, error) {
	return &TrailManagerService{
		repoTM:    rTM,
		repoT:     rT,
		repoS:     rS,
		repoZ:     rZ,
		publisher: publisher,
	}, nil
}

func (t *TrailManagerService) CreateTrailManager(wId uuid.UUID) (uuid.UUID, error) {
	tm, err := domain.NewTrailManager(wId)
	if err != nil {
		return uuid.Nil, err
	}
	t.repoTM.AddTrailManagerIntance(tm)
	return tm.TrailManagerID, nil
}

func (t *TrailManagerService) CloseTrailManager(wId uuid.UUID) error {
	return t.repoTM.DeleteTrailManagerInstance(wId)
}

func (t *TrailManagerService) GetTrailManagerByWorkoutId(id uuid.UUID) (*domain.TrailManager, error) {
	tm, err := t.repoTM.GetByWorkoutId(id)
	return tm, err
}

func (t *TrailManagerService) CreateTrail(tid uuid.UUID, name string, zId uuid.UUID, startLatitude float64, startLongitude float64, endLatitude float64, endLongitude float64) (uuid.UUID, error) {
	res, err := t.repoT.CreateTrail(tid, name, zId, startLatitude, startLongitude, endLatitude, endLongitude)
	if err != nil {
		return uuid.Nil, err
	}
	return res, nil
}

func (t *TrailManagerService) UpdateTrail(tid uuid.UUID, name string, zId uuid.UUID, startLatitude float64, startLongitude float64, endLatitude float64, endLongitude float64) error {
	err := t.repoT.UpdateTrailByID(tid, name, zId, startLatitude, startLongitude, endLatitude, endLongitude)
	if err != nil {
		return err
	}
	return nil
}

func (t *TrailManagerService) DeleteTrail(tId uuid.UUID) error {
	err := t.repoT.DeleteTrailByID(tId)
	log.Debug("deelte trail err", zap.Error(err))
	if err != nil {
		return err
	}
	return nil
}

func (t *TrailManagerService) DisconnectTrailManager(wId uuid.UUID) error {
	err := t.repoTM.DeleteTrailManagerInstance(wId)
	return err
}

func (t *TrailManagerService) GetTrailByID(id uuid.UUID) (*domain.Trail, error) {
	trail, err := t.repoT.GetTrailByID(id)
	return trail, err
}

func (t *TrailManagerService) CheckTrail(id uuid.UUID) error {
	trail, err := t.repoT.GetTrailByID(id)
	if err != nil || trail.TrailID != id {
		return err
	}
	return nil
}

func (t *TrailManagerService) GetCurrentLocation(wId uuid.UUID) (float64, float64, error) {
	tmInstance, err := t.repoTM.GetByWorkoutId(wId)
	if err != nil {
		return math.MaxFloat64, math.MaxFloat64, err
	}
	return tmInstance.CurrentLongitude, tmInstance.CurrentLatitude, nil
}

func (t *TrailManagerService) GetClosestTrail(zId uuid.UUID, currentLongitude float64, currentLatitude float64) (uuid.UUID, error) {

	trails, err := t.repoT.ListTrailsByZoneId(zId)
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

	// If a closest trail is found, update the TrailManager
	if closestTrail != nil {

		return closestTrail.TrailID, nil
	}

	return uuid.Nil, nil // Or return an appropriate error if necessary
}

func (t *TrailManagerService) CreateShelter(id uuid.UUID, name string, tId uuid.UUID, availability bool, lat, long float64) (uuid.UUID, error) {
	sId, err := t.repoS.CreateShelter(id, name, tId, availability, lat, long)
	if err != nil {
		return uuid.Nil, err
	} else {
		return sId, nil
	}
}

func (t *TrailManagerService) UpdateShelter(id uuid.UUID, name string, tId uuid.UUID, availability bool, lat, long float64) error {
	return t.repoS.UpdateShelterByID(id, tId, name, availability, lat, long)
}

func (t *TrailManagerService) DeleteShelter(id uuid.UUID) error {
	return t.repoS.DeleteShelterByID(id)
}

func (t *TrailManagerService) GetShelterByID(id uuid.UUID) (*domain.Shelter, error) {
	shelter, err := t.repoS.GetShelterByID(id)
	return shelter, err
}

func (t *TrailManagerService) CheckShelter(id uuid.UUID) error {
	s, err := t.repoS.GetShelterByID(id)
	if err != nil || s.ShelterID != id {
		return err
	}
	return nil
}

func (t *TrailManagerService) CreateZone(zName string) (uuid.UUID, error) {
	zId, err := t.repoZ.CreateZone(zName)
	if err != nil {
		return uuid.Nil, err
	}
	return zId, nil
}

func (t *TrailManagerService) CheckZone(zId uuid.UUID) error {
	z, err := t.repoZ.GetZoneByID(zId)
	if err != nil || z.ZoneID != zId {
		return err
	}
	return nil
}

func (t *TrailManagerService) UpdateZone(zId uuid.UUID, zName string) error {
	err := t.repoZ.UpdateZone(zId, zName)
	if err != nil {
		return err
	}
	return nil
}

func (t *TrailManagerService) DeleteZone(zId uuid.UUID) error {

	err := t.CheckZone(zId)
	if err != nil {
		return nil
	}

	err = t.repoZ.DeleteZone(zId)
	if err != nil {
		return err
	}
	return nil
}

func (t *TrailManagerService) UpdateCurrentLocation(wId uuid.UUID, latitude float64, longitude float64, time time.Time) error {
	tmInstance, err := t.repoTM.GetByWorkoutId(wId)
	if err != nil {
		log.Debug("No such trail manager instance found will create on request", zap.Error(err))
		_, _ = t.CreateTrailManager(wId)
		tmInstance, _ = t.repoTM.GetByWorkoutId(wId)
	}
	tmInstance.CurrentLatitude = latitude
	tmInstance.CurrentLongitude = longitude
	tmInstance.CurrentWorkoutID = wId
	tmInstance.CurrentTime = time
	t.repoTM.Update(*tmInstance)

	// Now push the shelter data data to the queue to the workout
	shelterId, distance, availability, _, err := t.GetClosestShelter(wId)
	if err != nil {
		log.Error("error when getting cloest shelter info", zap.Error(err))
		return err
	}
	closestShelter, _ := t.GetShelterByID(shelterId)
	err = t.publisher.PublishShelterInfo(shelterId, closestShelter.ShelterName, availability, distance)

	if err != nil {
		log.Error("error when publishing shelter info", zap.Error(err))
	}
	log.Debug("publishing shelter data to workout thru queue")

	return nil
}

func (t *TrailManagerService) GetClosestShelter(wId uuid.UUID) (uuid.UUID, float64, bool, time.Time, error) {

	shelters, err := t.repoS.List()
	if err != nil {
		return uuid.Nil, math.MaxFloat64, false, time.Now(), err
	}
	tm, err1 := t.repoTM.GetByWorkoutId(wId)
	if err1 != nil {
		return uuid.Nil, math.MaxFloat64, false, time.Now(), err1
	}

	if tm == nil {
		// if it does not exist, create a new instance
		tmNew, _ := domain.NewTrailManager(wId)
		t.repoTM.AddTrailManagerIntance(tmNew)
		tm, _ = t.repoTM.GetByWorkoutId(wId)
	}
	var closestShelter *domain.Shelter
	minDistance := math.MaxFloat64 // Initialize with the maximum float value

	for _, shelter := range shelters {

		distance := 0.0
		point1 := haversine.Coord{Lat: tm.CurrentLatitude, Lon: tm.CurrentLongitude}
		point2 := haversine.Coord{Lat: shelter.Latitude, Lon: shelter.Longitude}
		_, distance = haversine.Distance(point1, point2)

		if distance < minDistance {
			minDistance = distance
			closestShelter = shelter
		}
	}

	// If a closest trail is found, update the TrailManager
	if closestShelter != nil {

		return closestShelter.ShelterID, minDistance, closestShelter.ShelterAvailability, tm.CurrentTime, nil
	}

	return uuid.Nil, math.MaxFloat64, false, time.Now(), nil // Or return an appropriate error if necessary
}

func (t *TrailManagerService) GetClosestShelterInfo(latitude float64, longitude float64) (uuid.UUID, float64, error) {
	shelters, err := t.repoS.List()
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
