package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

/*
LS-TODO: Remove json parsing information & remove un-needed fields
LS-TODO Move the entities like HRM and Geo to separate structs, for example:
type HRM struct {

}
type Geo struct {}
*/
type Peripheral struct {
	PlayerId        uuid.UUID
	WorkoutId       uuid.UUID
	HRMId           uuid.UUID
	HRate           int
	HRateTime       time.Time
	HRMStatus       bool
	CreatedAt       time.Time
	LocationTime    time.Time
	GeoId           uuid.UUID
	GeoStatus       bool
	GeoBrodacasting bool
	Longitude       float64
	Latitude        float64
	HRateCount      int
	AverageHRate    int
	LiveData        bool
}

func (p *Peripheral) GetAverageHRate() LastHR {
	var tempHRDTO LastHR
	tempHRDTO.HRMID = p.HRMId
	tempHRDTO.TimeOfLocation = p.HRateTime
	tempHRDTO.HeartRate = p.AverageHRate
	return tempHRDTO
}

func (p *Peripheral) GetHRate() LastHR {
	var tempHRDTO LastHR
	tempHRDTO.HRMID = p.HRMId
	tempHRDTO.TimeOfLocation = p.HRateTime
	tempHRDTO.HeartRate = p.HRate
	return tempHRDTO
}

func (p *Peripheral) SetHRate(reading int) {
	if p.HRMStatus {
		p.AverageHRate = (p.HRate*p.HRateCount + reading) * 1.0 / (1 + p.HRateCount)
		p.HRateCount += 1
		p.HRate = reading
		p.HRateTime = time.Now()
	}
}

func (p *Peripheral) SetHRMStatus(code bool) {
	if code {
		p.HRMStatus = true
	} else {
		p.HRMStatus = false
	}
	fmt.Println(p.HRMStatus)
}

// function for getting the reading of longitude and lattide
func (p *Peripheral) SetLocation(longitude float64, latitude float64) {
	if p.GeoStatus {
		p.LocationTime = time.Now()
		p.Longitude = longitude
		p.Latitude = latitude
	}
}

// NOTES: ONLY read location if the peripheral status is on, otherwise it is off, so
func (p *Peripheral) GetGeoLocation() LastLocation {
	var tempLocationDTO LastLocation
	tempLocationDTO.TimeOfLocation = p.LocationTime
	tempLocationDTO.Longitude = p.Longitude
	tempLocationDTO.Latitude = p.Latitude
	tempLocationDTO.WorkoutID = p.WorkoutId
	return tempLocationDTO

}

// LS-TODO: The function should take the values as parameters. See User.go for example
// NewPeripheral is a factory to create a new Peripheral aggregate
func NewPeripheral(p Peripheral) (Peripheral, error) {

	// LS-TODO: Add validation for different fields.
	// Create a hrm object and initialize all the values to avoid nil pointer exceptions
	pN := Peripheral{
		PeripheralId:    uuid.New(),
		PlayerId:        p.PlayerId,
		HRMId:           p.HRMId,
		WorkoutId:       p.WorkoutId,
		GeoId:           uuid.New(),
		CreatedAt:       time.Now(),
		HRMStatus:       p.HRMStatus,
		GeoStatus:       true,
		LiveData:        p.LiveData,
		GeoBrodacasting: p.GeoBrodacasting,
	}
	return pN, nil
}
