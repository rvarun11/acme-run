package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// UserServiceClientMock is a mock type for UserServiceClientImpl
type UserServiceClientMock struct {
	mock.Mock
}

// NewUserServiceClientMock creates a new instance of UserServiceClientMock.
func NewUserServiceClientMock() *UserServiceClientMock {
	return &UserServiceClientMock{}
}

// GetProfileOfUser provides a mock function with given fields
func (m *UserServiceClientMock) GetProfileOfUser(playerID uuid.UUID) (string, error) {
	args := m.Called(playerID)
	return args.String(0), args.Error(1)
}

// GetHardcoreModeOfUser provides a mock function with given fields
func (m *UserServiceClientMock) GetHardcoreModeOfUser(playerID uuid.UUID) (bool, error) {
	args := m.Called(playerID)
	return args.Bool(0), args.Error(1)
}

// GetUserAge provides a mock function to get the age of a user
func (m *UserServiceClientMock) GetUserAge(playerID uuid.UUID) (uint8, error) {
	args := m.Called(playerID)
	return uint8(args.Int(0)), args.Error(1)
}
