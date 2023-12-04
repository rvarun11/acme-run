package amqp

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type ShelterDistancePublisherMock struct {
	mock.Mock
}

func NewShelterDistancePublisherMock() *ShelterDistancePublisherMock {
	return &ShelterDistancePublisherMock{}
}

func (m *ShelterDistancePublisherMock) PublishShelterDistance(wId uuid.UUID, sId uuid.UUID, name string, availability bool, distance float64) error {
	args := m.Called(wId, sId, name, availability, distance)
	return args.Error(0)
}
