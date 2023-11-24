package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

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

func (p *Peripheral) GetAverageHRate() (uuid.UUID, time.Time, int) {
	return p.HRMId, p.HRMDev.HRateTime, p.HRMDev.AverageHRate
}

func (p *Peripheral) GetHRate() (uuid.UUID, time.Time, int) {
	return p.HRMId, p.HRMDev.HRateTime, p.HRMDev.HRate
}

func (p *Peripheral) SetHRate(reading int) {
	if p.HRMDev.HRMStatus {
		p.HRMDev.AverageHRate = (p.HRMDev.AverageHRate*p.HRMDev.HRateCount + reading) * 1.0 / (1 + p.HRMDev.HRateCount)
		fmt.Println("count %d avg %d current %d", p.HRMDev.HRateCount, p.HRMDev.AverageHRate, reading)
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
func (p *Peripheral) GetGeoLocation() (time.Time, float64, float64, uuid.UUID) {

	return p.GeoDev.LocationTime, p.GeoDev.Longitude, p.GeoDev.Latitude, p.WorkoutId

}

func NewPeripheral(pId uuid.UUID, hId uuid.UUID, wId uuid.UUID, hStatus bool, liveStatus bool) (Peripheral, error) {

	if pId == uuid.Nil {
		return Peripheral{}, errors.New("player ID cannot be empty")
	}
	if hId == uuid.Nil {
		return Peripheral{}, errors.New("HRM ID cannot be empty")
	}
	if wId == uuid.Nil {
		return Peripheral{}, errors.New("workout ID cannot be empty")
	}

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
