package services

import (
	"math"

	"github.com/CAS735-F23/macrun-teamvsl/trail/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/trail/internal/core/ports"
	"github.com/google/uuid"
)

type TrailManagerService struct {
	repoT  ports.TrailRepository
	repoS  ports.ShelterRepository
	repoTM ports.TrailManagerRepository
}

// Factory for creating a new TrailManager

func NewTrailManagerService(rTM ports.TrailManagerRepository, rT ports.TrailRepository, rS ports.ShelterRepository) (*TrailManagerService, error) {
	return &TrailManagerService{
		repoTM: rTM,
		repoT:  rT,
		repoS:  rS,
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

func (t *TrailManagerService) CreateTrail(tid uuid.UUID, name string, startLatitude float64, startLongitude float64, endLatitude float64, endLongitude float64, shelterId uuid.UUID) (uuid.UUID, error) {
	res, err := t.repoT.CreateTrail(tid, name, startLatitude, startLongitude, endLatitude, endLongitude, shelterId)
	if err != nil {
		return uuid.Nil, err
	}
	return res, nil
}

func (t *TrailManagerService) DisconnectTrailManager(wId uuid.UUID) error {
	err := t.repoTM.DeleteTrailManagerInstance(wId)
	return err
}

func (t *TrailManagerService) SetCurrentLocation(wId uuid.UUID, longitude float64, latitude float64) error {
	tmInstance, err := t.repoTM.GetByWorkoutId(wId)
	if err != nil {
		return err
	}
	tmInstance.CurrentLongitude = longitude
	tmInstance.CurrentLatitude = latitude
	t.repoTM.Update(*tmInstance)
	return nil
}

func (t *TrailManagerService) GetShelter(id uuid.UUID) (*domain.Shelter, error) {
	shelter, err := t.repoS.GetShelterByID(id)
	return shelter, err
}

func (t *TrailManagerService) GetTrail(id uuid.UUID) (*domain.Trail, error) {
	trail, err := t.repoT.GetTrailByID(id)
	return trail, err
}

func (t *TrailManagerService) calculateDistance(Longitude1 float64, Latitude1, Longitude2 float64, Latitude2 float64) (float64, error) {
	x := Longitude2 - Longitude1
	y := Latitude2 - Latitude1
	return math.Sqrt(x*x + y*y), nil
}

func (t *TrailManagerService) SelectTrail(wId uuid.UUID, tId uuid.UUID, option string) error {
	tmInstance, err := t.repoTM.GetByWorkoutId(wId)
	if err != nil {
		return err
	}
	if option == "bind" {
		tmInstance.CurrentTrailID = tId
	} else if option == "unbind" {
		tmInstance.CurrentTrailID = uuid.Nil
	} else {
		return ports.ErrorTrailManagerInvalidArgs
	}
	t.repoTM.Update(*tmInstance)
	return nil
}

// function for compute the distance between current geo reading to the cloest shelter
func (t *TrailManagerService) GetShelterDistance(wId uuid.UUID, tId uuid.UUID, sId uuid.UUID) (float64, error) {

	tmInstance, err := t.repoTM.GetByWorkoutId(wId)
	if err != nil {
		return math.MaxFloat64, err
	}
	// Retrieve the details of the closest shelter from the repository.
	shelter, err := t.repoS.GetShelterByID(tmInstance.CurrentShelterID)
	if err != nil {
		// Handle the error if the shelter is not found.
		return math.MaxFloat64, err
	}

	// Convert latitude and longitude from degrees to radians.
	lon1 := tmInstance.CurrentLongitude
	lat1 := tmInstance.CurrentLatitude
	lon2 := shelter.Longitude
	lat2 := shelter.Latitude

	x := lon2 - lon1
	y := lat2 - lat1
	return math.Sqrt(x*x + y*y), nil

}

func (t *TrailManagerService) GetTrailDistance(wId uuid.UUID, tId uuid.UUID, sId uuid.UUID) (float64, error) {

	tmInstance, err := t.repoTM.GetByWorkoutId(wId)
	if err != nil {
		return math.MaxFloat64, err
	}
	// Retrieve the details of the closest shelter from the repository.
	trail, err := t.repoT.GetTrailByID(tId)
	if err != nil {
		// Handle the error if the shelter is not found.
		return math.MaxFloat64, err
	}

	// Convert latitude and longitude from degrees to radians.
	lon1 := tmInstance.CurrentLongitude
	lat1 := tmInstance.CurrentLatitude
	lon2 := trail.StartLongitude
	lat2 := trail.StartLatitude

	x := lon2 - lon1
	y := lat2 - lat1
	return math.Sqrt(x*x + y*y), nil

}

func (t *TrailManagerService) GetClosestShelter(currentLongitude, currentLatitude float64) (uuid.UUID, error) {
	shelters, err := t.repoS.List()
	if err != nil {
		return uuid.Nil, err // Handle the error, possibly no shelters available or DB error
	}

	var closestShelter *domain.Shelter
	minDistance := math.MaxFloat64 // Initialize with the maximum float value

	// Find the closest shelter
	for _, shelter := range shelters {
		distance, _ := t.calculateDistance(currentLongitude, currentLatitude, shelter.Longitude, shelter.Latitude)
		if distance < minDistance {
			minDistance = distance
			closestShelter = shelter
		}
	}

	if closestShelter != nil {
		return closestShelter.ShelterID, nil
	}
	return uuid.Nil, nil
}

func (s *TrailManagerService) GetClosestTrail(currentLongitude float64, currentLatitude float64) (uuid.UUID, error) {

	trails, err := s.repoT.List()
	if err != nil {
		return uuid.Nil, nil // Handle the error, possibly no trails available or DB error
	}

	var closestTrail *domain.Trail
	minDistance := math.MaxFloat64 // Initialize with the maximum float value

	for _, trail := range trails {
		distance, _ := s.calculateDistance(currentLongitude, currentLatitude, trail.StartLongitude, trail.StartLatitude)
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

// // This should be a method of TrailManagerService
// func (t *TrailManagerService) ListenForLocationUpdates(queueName string, wId uuid.UUID) {
// 	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq/")
// 	if err != nil {
// 		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
// 	}
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		log.Fatalf("Failed to open a channel: %s", err)
// 	}
// 	defer ch.Close()

// 	msgs, err := ch.Consume(
// 		queueName,
// 		"",
// 		true,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		log.Fatalf("Failed to register a consumer: %s", err)
// 	}

// 	forever := make(chan bool)

// 	go func() {
// 		for d := range msgs {
// 			var location dto.LocationDTO //
// 			if err := json.Unmarshal(d.Body, &location); err != nil {
// 				log.Printf("Error decoding JSON: %s", err)
// 				continue
// 			}
// 			tmInstance, err := t.repoTM.GetByWorkoutId(wId)
// 			if err != nil {

// 			}
// 			tmInstance.CurrentLongitude = location.Longitude
// 			tmInstance.CurrentLatitude = location.Latitude
// 			// Update TrailManager's current location using the repository method.
// 			// if err := t.repoTM.UpdateLocation(location.WorkoutID, location.Longitude, location.Latitude); err != nil {
// 			// 	log.Printf("Failed to update location: %s", err)
// 			// }
// 		}
// 	}()

// 	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
// 	<-forever
// }
