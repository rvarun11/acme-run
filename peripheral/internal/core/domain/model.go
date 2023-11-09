package domain

import (
	"errors"
	"time"
	"github.com/google/uuid"
)

var (
	ErrInvalidEmail = errors.New("a customer has to have a valid email")
)

type Peripheral struct {
	WorkoutId uuid.UUID `json:"workout_id"`
	HRMId     uuid.UUID `json:"hrm_id"`
	HRate     string    `json:"heart_rate"`
	Status string `json:"peripheral_status"`
	CreatedAt time.Time `json:"created_at"`
	LocationTime time.time `json:"locationTime"`
	Longitude  float64       `json:"longitude"`
	Latitude   float64   `json:"latitude"`
}

// Getters and Setters for HRM
func (p *Peripheral) GetID() uuid.UUID {
	return p.HRMId
}

func (p *Peripheral) SetID(id uuid.UUID) {
	p.HRMId = id
}

func (p *Peripheral) getHRate()(time.time, string) string {
	return time.Now(),p.HRate
}

func (p *Peripheral) readHRate() {
	p.HRate = "100"
}

func (p *Peripheral) SetStatus(code int) {
	if code == 1{
		p.Status = "on"
	}else {
		p.Status = "off"
	}
}

func (p *Peripheral) getStatus() string {
	return p.Status
}

// NOTES: ONLY read location if the peripheral status is on, otherwise it is off, so
func(p *Peripheral) getLocation() (time.time, float64, float64)
{
	p.LocationTime = time.Now()
	// TODO: implement a location random generator
	p.Longitude = 10.0
	p.Latitude = 20.0
	return p.LocationTime, p.Longitude, p.Latitude
}

// NewPlayer is a factory to create a new Player aggregate

func NewPeripheral(p Peripheral) (Peripheral, error) {

	// Create a hrm object and initialize all the values to avoid nil pointer exceptions
	pN := Peripheral{
		HRMId:     p.HRMId,
		CreatedAt: time.Now(),
		Status: "on"
	}
	return pN, nil
}
