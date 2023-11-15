// Package entities holds all the entities that are shared across all subdomains
package domain

import (
	"errors"
	"time"
    "database/sql"
    "log"
	"math"
    _ "github.com/lib/pq"
	"github.com/google/uuid"
)

var (
	ErrInvalidTrail                 = errors.New("no trail_id matched")
	ErrInvalidTrailManager                 = errors.New("no trail_manager_id matched")
)

type Shelter struct {
	// ID of the shelter
	ShelterID uuid.UUID `json:"shelter_id"`
	// zone id of shelter
	ZoneID uuid.UUID `json:"zone_id"`
	// availability of shelter
	ShelterAvailability bool `json:"shelter_available"`
	// name of the shelter
	ShelterName strig `json:"shelter_name"`
	// longitude of the shelter
	Longitude float64 `json:"longitude"`
	// latitude of the shelter
	Latitude float64 `json:"latitude"`
}

type Trail struct {
	// id of the trail
	TrailID uuid.UUID `json:"trail_id"`
	// name of the trail
	TrailName string `json:"trail_name"`
	// zone id of the trail
	ZoneID uui.UUID `json:"zone_id"`
	// start longitude
	StartLongitude float64 `json:"start_longitude"`
	// start latitude
	StartLatitude float64 `json:"start_latitude"`
	// end longitude
	EndLongitude float64 `json:"end_lontitude"`
	// end latitude
	EndLatitude float64 `json:"end_latitude"`
	// id of the cloest shelter
	CloestShelterId uuid.UUID `json:"cloest_shelter_id"`
	// created time
	CreatedAt time.time `json:"created_at"`
}




func (t *Trail) GetID() uuid.UUID {
	return t.TrailID
}

func (t *Trail) SetID(id uuid.UUID) {
	t.TrailID = id
}

func (t *Trail) GetCloestTrailID() uuid.UUID {
	return t.CloestShelterID
}

func (t *Trail) CheckTrailShelterAvailable bool{
	if t.CloestShelterID == uuid.nil{
		return false
	}
}

func newTrail(tId uuid.UUID, tName string, zId uuid.UUID, sLong float64, sLat float64, eLong float64, eLat float64, sId uuid.UUID) {
	if tId == uuid.Nil {
		return Trail{}, ErrInvalidTrail
	}

	return Trail{
		TrailID:         tId,
		TrailName:       tName,
		ZoneID: zId
		StartLongitude: sLong,
		StartLatitude:   sLat,
		EndLongitude:   eLong,
		EndLatitude:     eLat,
		ShelterID:       sId,
		CreatedAt:       time.Now(),
	}, nil
}

// Workout is a entity that represents a workout in all Domains
type TrailManager struct {
	// ID is the identifier of the Entity, the ID is shared for all sub domains
	TrailManagerID uuid.UUID `json:"trail_manager_id"`
	// record of current workout id it is tracking
	CurrentWorkoutID uuid.UUID `json:"current_workout_id"`
	// trailId is the id of the trail player is on
	CurrentTrailID uuid.UUID `json:"current_trail_id"`
	// record of current longitude
	CurrentLongitude float64 `json:"current_longitude"`
	// record of current latitude
	CurrentLatitude float64 `json:"current_latitude"`
	// record of current time
	CurrentTime time.time `json:"current_time"`
	// current shelter that is the cloest
	CloestShelterID uuid.UUID `json:"cloest_shelter_id"`
	// CreatedAt is the time when the trail manager was started
	CreatedAt time.Time `json:"created_at"`
}

func (tM *TrailManager) setCurrentTrail(tId uuid.UUID){
	tm.CurrentTrailID = tId
}



// function to get shelter by ID
func (t *Trail) GetShelterByID(shelterID uuid.UUID) (*Shelter, error) {
	if shelterID == uuid.Nil {
		return nil, ErrInvalidShelter
	}
	// TODO
	db, err := sql.Open("postgres", "your-postgres-connection-string-for-shelters")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var shelter Shelter
	query := `SELECT shelter_id, shelter_available, shelter_name, longitude, latitude FROM shelters WHERE shelter_id = $1`
	row := db.QueryRow(query, shelterID)

	err = row.Scan(&shelter.ShelterID, &shelter.ShelterAvailability, &shelter.ShelterName, &shelter.Longitude, &shelter.Latitude)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidShelter
		}
		return nil, err
	}

	return &shelter, nil
}

// function for getting trail by ID
func (s *TrailManagerService) GetTrailByID(trailID uuid.UUID) (*Trail, error) {
	if trailID == uuid.Nil {
		return nil, ErrInvalidTrail
	}
	// TODO
	db, err := sql.Open("postgres", "your-postgres-connection-string")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var trail Trail
	query := `SELECT trail_id, trail_name,zone_id, start_longitude, start_latitude, end_longitude, end_latitude, closest_shelter_id, created_at FROM trails WHERE trail_id = $1`
	row := db.QueryRow(query, trailID)

	err = row.Scan(&trail.TrailID, &trail.TrailName, &trail.StartLongitude, &trail.StartLatitude, &trail.EndLongitude, &trail.EndLatitude, &trail.ClosestShelterID, &trail.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidTrail
		}
		return nil, err
	}

	return &trail, nil
}

func (tM *TrailManagerService) calculateDistance(lon1, lat1, lon2, lat2 float64) (float64,error){ 

	// Use Pythagorean Theorem on an equirectangular projection
	x := lon2 - lon1
	y := lat2 - lat1
	return math.Sqrt(x*x + y*y) , nil
}

// function for compute the distance between current geo reading to the cloest shelter
func (tM *TrailManager) getCurrentShelterDistance(db *sql.DB, shelterID uuid.UUID) (float64, error) {
	if tm.CurrentLatitude == 0 || tm.CurrentLongitude == 0 {
		return 0, errors.New("current location of trail manager is not set")
	}

	shelter, err := getShelterByID(db, shelterID)
	if err != nil {
		return 0, err
	}

	// distance := calculateDistance(tm.CurrentLatitude, tm.CurrentLongitude, shelter.Latitude, shelter.Longitude)

	lon1 := tM.CurrentLongitude
	lat1 := tM.CurrentLatitude
	lon2 := shelter.Longitude
	lat2 := shelter.Latitude

	// Use Pythagorean Theorem on an equirectangular projection
	x := lon2 - lon1
	y := lat2 - lat1
	return math.Sqrt(x*x + y*y) , nil


}

// failOnError is a helper function to log any errors and fail fast.
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type LocationDTO struct {
	WorkoutId    UUID.UUID `json:"workout_id"`
	LocationTime time.Time `json:"locationTime"`
	Longitude    float64   `json:"longitude"`
	Latitude     float64   `json:"latitude"`
}

// function for grab the current location of the 
// func (tM *TrailManager) getCurrentLocation(){
// 	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq/")
// 	failOnError(err, "Failed to connect to RabbitMQ")
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	failOnError(err, "Failed to open a channel")
// 	defer ch.Close()

// 	q, err := ch.QueueDeclare(
// 		"HR-Queue-002", // Queue name must match the one used by the PeripheralService
// 		false,          // durable
// 		false,          // delete when unused
// 		false,          // exclusive
// 		false,          // no-wait
// 		nil,            // arguments
// 	)
// 	failOnError(err, "Failed to declare a queue")

// 	msgs, err := ch.Consume(
// 		q.Name, // queue
// 		"",     // consumer
// 		false,  // auto-ack
// 		false,  // exclusive
// 		false,  // no-local
// 		false,  // no-wait
// 		nil,    // args
// 	)
// 	failOnError(err, "Failed to register a consumer")

// 	var latestLocation LocationDTO
// 	var latestDeliveryTag uint64

// 	for d := range msgs {
// 		var location LocationDTO
// 		err := json.Unmarshal(d.Body, &location)
// 		if err != nil {
// 			log.Printf("Error decoding JSON: %s", err)
// 			continue
// 		}

// 		// Assuming LocationTime is used to determine the latest message.
// 		if latestLocation.LocationTime.Before(location.LocationTime) {
// 			latestLocation = location
// 			latestDeliveryTag = d.DeliveryTag
// 		}

// 		// Acknowledge the previous latest message
// 		if latestDeliveryTag != 0 {
// 			ch.Ack(latestDeliveryTag, false)
// 		}
// 	}

// 	tM.CurrentLongitude = location.Longitude
// 	tM.CurrentLatitude = location.Latitude
// 	tM.CurrentTime = location.LocationTime
// 	tM.CurrentWorkID = location.WorkoutId
// }

func NewTrailManager() (TrailManager, error) {

	return TrailManager{
		TrailManagerID: uuid.New(),
		CurrentWorkoutID: uuid.nil,
		CurrentTrailID: uuid.nil,
		CurrentLongitude: 0.0,
		CurrentLatitude: 0.0,
		CurrentTime: time.Now(),
		CloestShelterID: uuid.nil
		CreatedAt time.Now(),

	}, nil
}

