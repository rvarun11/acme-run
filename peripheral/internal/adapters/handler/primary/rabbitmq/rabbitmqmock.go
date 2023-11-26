package rabbitmqhandler

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// RabbitMQHandlerMock is a mock type for the RabbitMQHandler type
type RabbitMQHandlerMock struct {
	mock.Mock
}

// NewRabbitMQHandlerMock creates a new instance of RabbitMQHandlerMock
func NewRabbitMQHandlerMock() *RabbitMQHandlerMock {
	return &RabbitMQHandlerMock{}
}

// SendLastLocation is a mock method that simulates sending location data to a queue
func (r *RabbitMQHandlerMock) SendLastLocation(wId uuid.UUID, latitude float64, longitude float64, time time.Time, toTrail bool) error {
	args := r.Called(wId, latitude, longitude, time)
	return args.Error(0)
}
