package domain

import (
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/adapters/dto"
	"github.com/google/uuid"
)

/*
LS-TODO: Remove json parsing information & remove un-needed fields
LS-TODO Move the entities like HRM and Geo to separate structs, for example:
type HRM struct {

}
type Geo struct {}
*/
type HRMData struct {
	HRate        int
	HRateTime    time.Time
	HRMStatus    bool
	HRateCount   int
	AverageHRate int
}

type GeoData struct {
	LocationTime time.Time
	GeoStatus    bool
	Longitude    float64
	Latitude     float64
}

type Peripheral struct {
	PlayerId   uuid.UUID
	WorkoutId  uuid.UUID
	HRMId      uuid.UUID
	HRMDev     HRMData
	GeoDev     GeoData
	CreatedAt  time.Time
	LiveStatus bool
}

func (p *Peripheral) GetAverageHRate() dto.LastHR {
	var tempHRDTO dto.LastHR
	tempHRDTO.HRMID = p.HRMId
	tempHRDTO.TimeOfLocation = p.HRMDev.HRateTime
	tempHRDTO.HeartRate = p.HRMDev.AverageHRate
	return tempHRDTO
}

func (p *Peripheral) GetHRate() dto.LastHR {
	var tempHRDTO dto.LastHR
	tempHRDTO.HRMID = p.HRMId
	tempHRDTO.TimeOfLocation = p.HRMDev.HRateTime
	tempHRDTO.HeartRate = p.HRMDev.HRate
	return tempHRDTO
}

func (p *Peripheral) SetHRate(reading int) {
	if p.HRMDev.HRMStatus {
		p.HRMDev.AverageHRate = (p.HRMDev.HRate*p.HRMDev.HRateCount + reading) * 1.0 / (1 + p.HRMDev.HRateCount)
		p.HRMDev.HRateCount += 1
		p.HRMDev.HRate = reading
		p.HRMDev.HRateTime = time.Now()
	}
}

// function for getting the reading of longitude and lattide
func (p *Peripheral) SetLocation(longitude float64, latitude float64) {
	if p.GeoDev.GeoStatus {
		p.GeoDev.LocationTime = time.Now()
		p.GeoDev.Longitude = longitude
		p.GeoDev.Latitude = latitude
	}
}

// NOTES: ONLY read location if the peripheral status is on, otherwise it is off, so
func (p *Peripheral) GetGeoLocation() dto.LastLocation {
	var tempLocationDTO dto.LastLocation
	tempLocationDTO.TimeOfLocation = p.GeoDev.LocationTime
	tempLocationDTO.Longitude = p.GeoDev.Longitude
	tempLocationDTO.Latitude = p.GeoDev.Latitude
	tempLocationDTO.WorkoutID = p.WorkoutId
	return tempLocationDTO

}

func NewPeripheral(pId uuid.UUID, hId uuid.UUID, wId uuid.UUID, hStatus bool, liveStatus bool) (Peripheral, error) {

	// LS-TODO: Add validation for different fields.
	// Create a hrm object and initialize all the values to avoid nil pointer exceptions
	pN := Peripheral{
		PlayerId:   pId,
		HRMId:      hId,
		WorkoutId:  wId,
		CreatedAt:  time.Now(),
		HRMDev:     HRMData{},
		GeoDev:     GeoData{},
		LiveStatus: liveStatus,
	}
	pN.HRMDev.HRMStatus = hStatus
	pN.GeoDev.GeoStatus = true

	return pN, nil
}
