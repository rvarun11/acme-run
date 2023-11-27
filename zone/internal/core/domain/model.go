package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var (
	ErrInvalidTrail       = errors.New("no trail_id matched")
	ErrInvalidZoneManager = errors.New("no trail_manager_id matched")
)

type Shelter struct {
	// ID of the shelter
	ShelterID uuid.UUID
	// the trail it is attached to
	TrailID uuid.UUID
	// availability of shelter
	ShelterAvailability bool
	// name of the shelter
	ShelterName string
	// longitude of the shelter
	Longitude float64
	// latitude of the shelter
	Latitude float64
}

type Trail struct {
	// id of the trail
	TrailID uuid.UUID
	// name of the trail
	TrailName string
	// zone of the trail
	ZoneID uuid.UUID
	// start longitude
	StartLongitude float64
	// start latitude
	StartLatitude float64
	// end longitude
	EndLongitude float64
	// end latitude
	EndLatitude float64
	// created time
	CreatedAt time.Time
}

type Zone struct {
	// id the zone
	ZoneID uuid.UUID
	// name of the zone
	ZoneName string
}

// func (t *Trail) CheckTrailShelterAvailable() (bool, error) {
// 	if t.ShelterID == uuid.Nil {
// 		return false, nil
// 	} else {
// 		return true, nil
// 	}
// }

func newTrail(tId uuid.UUID, tName string, zId uuid.UUID, sLong float64, sLat float64, eLong float64, eLat float64) (Trail, error) {
	if tId == uuid.Nil {
		return Trail{}, ErrInvalidTrail
	}

	return Trail{
		TrailID:        tId,
		TrailName:      tName,
		ZoneID:         zId,
		StartLongitude: sLong,
		StartLatitude:  sLat,
		EndLongitude:   eLong,
		EndLatitude:    eLat,
		CreatedAt:      time.Now(),
	}, nil
}

func newZone(zId uuid.UUID, zName string) (Zone, error) {
	if zId == uuid.Nil {
		return Zone{}, ErrInvalidTrail
	}
	return Zone{
		ZoneID:   zId,
		ZoneName: zName,
	}, nil
}

func newShelter(sId uuid.UUID, tId uuid.UUID, availability bool, sName string, long float64, lat float64) (Shelter, error) {

	return Shelter{
		ShelterID:           sId,
		TrailID:             tId,
		ShelterAvailability: availability,
		ShelterName:         sName,
		Longitude:           long,
		Latitude:            lat,
	}, nil

}

// Workout is a entity that represents a workout in all Domains
type ZoneManager struct {
	// ID is the identifier of the Entity, the ID is shared for all sub domains
	ZoneManagerID uuid.UUID
	// record of current workout id it is tracking
	CurrentWorkoutID uuid.UUID
	// Current Zone ID
	ZoneID uuid.UUID
	// trailId is the id of the trail player is on
	CurrentTrailID uuid.UUID
	// record of current longitude
	CurrentLongitude float64
	// record of current latitude
	CurrentLatitude float64
	// record of current time
	CurrentTime time.Time
	// current shelter that is the cloest
	// CreatedAt is the time when the trail manager was started
	CreatedAt time.Time
}

func NewZoneManager(wId uuid.UUID) (ZoneManager, error) {

	return ZoneManager{
		ZoneManagerID:    uuid.New(),
		CurrentWorkoutID: wId,
		ZoneID:           uuid.Nil,
		CurrentTrailID:   uuid.Nil,
		CurrentLongitude: 0.0,
		CurrentLatitude:  0.0,
		CurrentTime:      time.Now(),
		CreatedAt:        time.Now(),
	}, nil
}
