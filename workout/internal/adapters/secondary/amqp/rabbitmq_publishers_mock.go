package amqphandler

import (
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/domain"
	logger "github.com/CAS735-F23/macrun-teamvsl/workout/log"
	"go.uber.org/zap"
)

// MockPublisher is a mock implementation of the Publisher interface
type MockPublisher struct {
	// Add fields to store information about calls to the methods, if necessary
	PublishedWorkouts []*domain.Workout
}

// NewMockPublisher creates a new instance of MockPublisher
func NewMockPublisher() *MockPublisher {
	return &MockPublisher{
		PublishedWorkouts: make([]*domain.Workout, 0),
	}
}

// PublishWorkoutStats mocks the PublishWorkoutStats method of Publisher
func (m *MockPublisher) PublishWorkoutStats(workoutStats *domain.Workout) error {
	// In the mock, we just store the workoutStats for verification in tests

	var challengeStatsDTO = challengeStatsDTO{
		PlayerID:        workoutStats.PlayerID,
		WorkoutEnd:      workoutStats.EndedAt,
		EnemiesFought:   workoutStats.Fights,
		EnemiesEscaped:  workoutStats.Escapes,
		DistanceCovered: workoutStats.DistanceCovered,
	}
	m.PublishedWorkouts = append(m.PublishedWorkouts, workoutStats)
	logger.Debug("workout statistics published to challenge manager", zap.Any("stats", challengeStatsDTO))
	return nil // Return nil to simulate successful execution
}
