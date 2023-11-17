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
	repoZ  ports.ZoneRepository
	repoTM ports.TrailManagerRepository
}

// Factory for creating a new TrailManager

func NewTrailManagerService(rTM ports.TrailManagerRepository, rT ports.TrailRepository, rS ports.ShelterRepository, rZ ports.ZoneRepository) (*TrailManagerService, error) {
	return &TrailManagerService{
		repoTM: rTM,
		repoT:  rT,
		repoS:  rS,
		repoZ:  rZ,
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

func (t *TrailManagerService) CheckTrailShelter(tId uuid.UUID) (bool, error) {
	trail, err := t.repoT.GetTrailByID(tId)
	if err != nil {
		return false, err
	}
	if trail.ShelterID == uuid.Nil {
		return false, nil
	} else {
		return true, nil
	}
}

func (t *TrailManagerService) GetCurrentLocation(wId uuid.UUID) (float64, float64, error) {
	tmInstance, err := t.repoTM.GetByWorkoutId(wId)
	if err != nil {
		return math.MaxFloat64, math.MaxFloat64, err
	}
	return tmInstance.CurrentLongitude, tmInstance.CurrentLatitude, nil
}

func (t *TrailManagerService) CalculateDistance(Longitude1 float64, Latitude1, Longitude2 float64, Latitude2 float64) (float64, error) {
	x := Longitude2 - Longitude1
	y := Latitude2 - Latitude1
	return math.Sqrt(x*x + y*y), nil
}

func (t *TrailManagerService) GetClosestTrail(zId uuid.UUID, currentLongitude float64, currentLatitude float64) (uuid.UUID, error) {

	trails, err := t.repoT.ListTrailsByZoneId(zId)
	if err != nil {
		return uuid.Nil, err // Handle the error, possibly no trails available or DB error
	}
	var closestTrail *domain.Trail
	minDistance := math.MaxFloat64 // Initialize with the maximum float value

	for _, trail := range trails {
		distance, _ := t.CalculateDistance(currentLongitude, currentLatitude, trail.StartLongitude, trail.StartLatitude)
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

func (t *TrailManagerService) CreateZone(zName string) (uuid.UUID, error) {
	zId, err := t.repoZ.CreateZone(zName)
	if err != nil {
		return uuid.Nil, err
	}
	return zId, nil
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
