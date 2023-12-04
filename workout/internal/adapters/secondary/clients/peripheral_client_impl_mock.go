package clients

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

// PeripheralClientMock is a mock type for PeripheralClientImpl
type PeripheralClientMock struct {
	mock.Mock
}

// NewPeripheralClientMock creates a new instance of PeripheralClientMock.
func NewPeripheralClientMock() *PeripheralClientMock {
	return &PeripheralClientMock{}
}

// BindPeripheralData provides a mock function with given fields
func (m *PeripheralClientMock) BindPeripheralData(trailID uuid.UUID, playerID uuid.UUID, workoutID uuid.UUID, hrmID uuid.UUID, HRMConnected bool, SendLiveLocationToTrailManager bool) error {
	args := m.Called(trailID, playerID, workoutID, hrmID, HRMConnected, SendLiveLocationToTrailManager)
	return args.Error(0)
}

// UnbindPeripheralData provides a mock function with given fields
func (m *PeripheralClientMock) UnbindPeripheralData(workoutID uuid.UUID) error {
	args := m.Called(workoutID)
	return args.Error(0)
}

// GetAverageHeartRateOfUser provides a mock function with given fields
func (m *PeripheralClientMock) GetAverageHeartRateOfUser(workoutID uuid.UUID) (uint8, error) {
	args := m.Called(workoutID)
	return args.Get(0).(uint8), args.Error(1)
}
