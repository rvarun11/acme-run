package clients

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type ZoneClientMock struct {
	mock.Mock
}

func NewZoneServiceClientMock() *ZoneClientMock {
	return &ZoneClientMock{}
}

func (z *ZoneClientMock) GetTrailLocation(trailID uuid.UUID) (float64, float64, float64, float64, error) {
	args := z.Called(trailID)
	return args.Get(0).(float64), args.Get(1).(float64), args.Get(2).(float64), args.Get(3).(float64), args.Error(4)
}
