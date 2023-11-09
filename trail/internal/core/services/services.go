package services

import (
	"fmt"
	"time"
	"math"
	"github.com/CAS735-F23/macrun-teamvsl/trail/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/trail/internal/core/ports"
	"github.com/google/uuid"
)


type TrailManagerService struct {
	tm *domain.TrailManager
	repoT ports.TrailRepository
	repoS ports.ShelterRepository
}

// Factory for creating a new TrailManager

func NewTrailManagerService(rT ports.TrailRepository, rS ports.ShelterRepository) *TrailManagerService{
	return &(TrailManagerService){
		tm: domain.NewTrailManager()
		repoT: rT,
		repoS, rS,
	}
}

// function for compute the distance between current geo reading to the cloest shelter
func (t *TrailManagerServices) getShelterDistance() (float64, error) {
	if t.tm.ClosestShelterID == uuid.Nil {
		// If the TrailManager has no closest shelter, return the maximum float value.
		return math.MaxFloat64, nil
	}

	// Retrieve the details of the closest shelter from the repository.
	shelter, err := t.repoS.GetShelterByID(t.tm.ClosestShelterID)
	if err != nil {
		// Handle the error if the shelter is not found.
		return math.MaxFloat64, err
	}

	// Convert latitude and longitude from degrees to radians.
	lon1 := t.tm.CurrentLongitude
	lat1 := t.tm.CurrentLatitude
	lon2 := shelter.Longitude
	lat2 := shelter.Latitude

	x := (lon2 - lon1R) * math.Cos((lat1+lat2)/2)
	y := lat2 - lat1
	return math.Sqrt(x*x + y*y) , nil

}

// TODO: DEPRECEATED
func (t *TrailManagerService) getClosestShelter(currentLongitude, currentLatitude float64) error {
	shelters, err := t.repoS.GetAllShelters()
	if err != nil {
		return err // Handle the error, possibly no shelters available or DB error
	}

	var closestShelter *domain.Shelter
	minDistance := math.MaxFloat64 // Initialize with the maximum float value

	// Find the closest shelter
	for _, shelter := range shelters {
		distance := t.calculateDistance(currentLongitude, currentLatitude, shelter.Longitude, shelter.Latitude)
		if distance < minDistance {
			minDistance = distance
			closestShelter = shelter
		}
	}

	// If a closest shelter is found, update the TrailManager
	if closestShelter != nil {
		t.tm.ClosestShelterID = closestShelter.ShelterID
		return nil
	}

	// If no shelter is close enough or there are no shelters, handle accordingly
	t.tm.ClosestShelterID = uuid.Nil
	return nil // Or return an appropriate error if necessary
}

func (t *TrailManagerService) getClosestShelterId(currentLongitude, currentLatitude float64) error {
	shelters, err := t.repoS.GetAllShelters()
	if err != nil {
		return err // Handle the error, possibly no shelters available or DB error
	}

	var closestShelter *domain.Shelter
	minDistance := math.MaxFloat64 // Initialize with the maximum float value

	// Find the closest shelter
	for _, shelter := range shelters {
		distance := t.calculateDistance(currentLongitude, currentLatitude, shelter.Longitude, shelter.Latitude)
		if distance < minDistance {
			minDistance = distance
			closestShelter = shelter
		}
	}

	if closestShelter != nil {
		
		return closestShelter.ShelterID
	}
	return uuid.Nil 
}



func (t *TrailManagerService) GetShelter(id uuid.UUID) (*domain.Shelter, error) {
	return t.repoS.GetShelterByID(id)
}

func (t *TrailManagerService) GetTrail(id uuid.UUID) (*domain.Trail, error) {
	return t.repoT.GetTrailByID(id)
}

func (t *TrailManagerService) sendDistance(){
	wId := t.tm.CurrentWorkID
	s_available := (t.tm.CloestShelterID == uuid.nil)
	dis := t.getShelterDistance()
	publishShelter(wId,s_available,dis)
}


func (s *TrailManagerService) getClosestTrailID(currentLongitude, currentLatitudefloat64) uuid.UUID {

	trails, err := t.repoT.GetAllTrails()
	if err != nil {
		return uuid.nil// Handle the error, possibly no trails available or DB error
	}

	var closestTrail *domain.Trail
	minDistance := math.MaxFloat64 // Initialize with the maximum float value

	for _, trail := range trails {
		distance := s.calculateDistance(currentLongitude, currentLatitude, trail.StartLongitude, trail.StartLatitude)
		if distance < minDistance {
			minDistance = distance
			closestTrail = trail
		}
	}

	// If a closest trail is found, update the TrailManager
	if closestTrail != nil {
		
		return uuid.nil
	}

	return closestTrail.TrailID // Or return an appropriate error if necessary	
}

