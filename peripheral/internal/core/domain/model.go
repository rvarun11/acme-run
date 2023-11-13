package domain

import (
	"errors"
	"time"
	"github.com/google/uuid"
)

var (
	
)

type Peripheral struct {
	PeripheralId uuid.UUID `json:"peripheral_id"`
	PlayerId uuid.UUID `json:"player_id"`
	WorkoutId uuid.UUID `json:"workout_id"`
	HRMId     uuid.UUID `json:"hrm_id"`
	HRate     int    `json:"heart_rate"`
	HRateTime time.time `json:"hrate_time"`
	HRMStatus bool `json:"hrm_status"`
	CreatedAt time.Time `json:"created_at"`
	LocationTime time.time `json:"locationTime"`
	GeoId uuid.UUID `json:"geo_id"`
	GeoStatus bool `json:"geo_status"`
	GeoBrodacasting bool `json:"geo_broadcasting"`
	Longitude  float64       `json:"longitude"`
	Latitude   float64   `json:"latitude"`
	HRateCount int `json:"heart_rate_count"`
	AverageHRate int `json:"average_heart_rate`
}

// Getters and Setters for HRM
func (p *Peripheral) GetHRMID() uuid.UUID {
	return p.HRMId
}

func (p *Peripheral) SetHRMID(id uuid.UUID) {
	p.HRMId = id
}

func (p *Peripheral) GetPeripheralID() uuid.UUID {
	return p.PeripheralId
}

func (p *Peripheral) GetGeoID(id uuid.UUID) {
	return p.GeoId
}

func (p *Peripheral) GetAverageHRate() LastHR {
	var tempHRDTO LastHR
	tempHRDTO.HRMID = p.HRMID
	tempHRDTO.TimeOfLocation = p.HRateTime
	tempHRDTO.HeartRate = p.AverageHRate
	return tempHRDTO
}

func (p *Peripheral) SetAverageHRate(rate int) {
	if p.HRMStatus == true {
		p.AverageHRate = ( p.HRate*p.HRateCount + rate )*1.0/(1+p.HRateCount)
		p.HRateCount += 1
		p.HRateTime = time.Now()
	}
}

func (p *Peripheral) SetHRMStatus(code bool) {
	if code == true{
		p.HRMStatus = true
	}else {
		p.HRMStatus = false
	}
}

// return the current status of the hrm
func (p *Peripheral) GetHRMStatus() bool {
	return p.HRMStatus
}

// return the current status of the hrm
func (p *Peripheral) GetGeoStatus() bool {
	return p.GeoStatus
}

func (p *Peripheral) SetGeoStatus(code bool) {
	if code == true{
		p.GeoStatus = true
	}else {
		p.GeoStatus = false
	}
}

// function for getting the reading of longitude and lattide
func(p *Peripheral) SetLocation(longitude float64, latitude float64) 
{
	if p.GeoStatus == true{
		p.LocationTime = time.Now()
		p.Longitude = longitude
		p.Latitude = latitude
	}
}

// NOTES: ONLY read location if the peripheral status is on, otherwise it is off, so
func(p *Peripheral) GetLocation() LastLocation
{
	var tempLocationDTO LastLocation
	tempLocationDTO.LocationTime = p.LocationTime
	tempLocationDTO.Longitude = p.Latitude
	tempLocationDTO.LocationTime = p.LocationTime
	tempLocationDTO.WorkoutID = p.WorkoutId
	return tempLocationDTO

}

// NewPlayer is a factory to create a new Player aggregate

func NewPeripheral(p Peripheral) (Peripheral, error) {

	// Create a hrm object and initialize all the values to avoid nil pointer exceptions
	pN := Peripheral{
		PeripheralId: uuid.New(),
		PlayerId: p.PlayerId,
		HRMId:     p.HRMId,
		WorkoutId: p.WorkoutId, 
		p.geo_id = uuid.New()
		CreatedAt: time.Now(),
		HRMStatus: p.HRMStatus,
		GeoBrodacasting: p.GeoBrodacasting	
	}
	return pN, nil
}
