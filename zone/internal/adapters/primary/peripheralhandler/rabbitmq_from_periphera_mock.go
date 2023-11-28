package peripheralhandler

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/mock"
)

type LocationSubscriberMock struct {
	mock.Mock
}

func (m *LocationSubscriberMock) CreateChannel(exchangeName, queueName, bindingKey, consumerTag string) (*amqp.Channel, error) {
	args := m.Called(exchangeName, queueName, bindingKey, consumerTag)
	return args.Get(0).(*amqp.Channel), args.Error(1)
}

func (m *LocationSubscriberMock) StartConsumer(workerPoolSize int, exchange, queueName, bindingKey, consumerTag string) error {
	args := m.Called(workerPoolSize, exchange, queueName, bindingKey, consumerTag)
	return args.Error(0)
}

func (m *LocationSubscriberMock) worker(deliveries <-chan amqp.Delivery) {
	m.Called(deliveries)
	// Implement the logic as needed or leave empty if not required for tests
}

type ZoneManagerAMQPHandlerMock struct {
	mock.Mock
}

func (m *ZoneManagerAMQPHandlerMock) InitAMQP() error {
	args := m.Called()
	return args.Error(0)
}
