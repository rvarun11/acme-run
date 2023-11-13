package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// Define types to hold request data for verification in tests
type BindRequestData struct {
	PlayerID                       uuid.UUID
	WorkoutID                      uuid.UUID
	HRMId                          uuid.UUID
	HRMConnected                   bool
	SendLiveLocationToTrailManager bool
}

type UnbindRequestData struct {
	WorkoutID uuid.UUID
}

// PeripheralDeviceClientMock is a mock type for PeripheralDeviceClientImpl
type PeripheralDeviceClientMock struct {
	mock.Mock
}

// NewPeripheralDeviceClientMock creates a new instance of PeripheralDeviceClientMock.
func NewPeripheralDeviceClientMock() *PeripheralDeviceClientMock {
	return &PeripheralDeviceClientMock{}
}

// BindPeripheralData provides a mock function with given fields
func (m *PeripheralDeviceClientMock) BindPeripheralData(playerID uuid.UUID, workoutID uuid.UUID, hrmID uuid.UUID, HRMConnected bool, SendLiveLocationToTrailManager bool) error {
	args := m.Called(playerID, workoutID, hrmID, HRMConnected, SendLiveLocationToTrailManager)
	return args.Error(0)
}

// UnbindPeripheralData provides a mock function with given fields
func (m *PeripheralDeviceClientMock) UnbindPeripheralData(workoutID uuid.UUID) error {
	args := m.Called(workoutID)
	return args.Error(0)
}

// GetAverageHeartRateOfUser provides a mock function with given fields
func (m *PeripheralDeviceClientMock) GetAverageHeartRateOfUser(workoutID uuid.UUID) (uint8, error) {
	args := m.Called(workoutID)
	return args.Get(0).(uint8), args.Error(1)
}
