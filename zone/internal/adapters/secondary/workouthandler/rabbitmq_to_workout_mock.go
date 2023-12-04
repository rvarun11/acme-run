package workouthandler

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type AMQPPublisherMock struct {
	mock.Mock
}

func NewAMQPPublisherMock() *AMQPPublisherMock {
	return &AMQPPublisherMock{}
}

func (m *AMQPPublisherMock) PublishShelterInfo(wId uuid.UUID, sId uuid.UUID, name string, availability bool, distance float64) error {
	args := m.Called(wId, sId, name, availability, distance)
	return args.Error(0)
}
