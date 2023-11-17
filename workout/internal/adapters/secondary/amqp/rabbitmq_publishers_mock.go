package amqphandler

import (
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/domain"
	logger "github.com/CAS735-F23/macrun-teamvsl/workout/log"
	"go.uber.org/zap"
)

// MockAMQPPublisher is a mock implementation of the AMQPPublisher interface
type MockAMQPPublisher struct {
	// Add fields to store information about calls to the methods, if necessary
	PublishedWorkouts []*domain.Workout
}

// NewMockAMQPPublisher creates a new instance of MockAMQPPublisher
func NewMockAMQPPublisher() *MockAMQPPublisher {
	return &MockAMQPPublisher{
		PublishedWorkouts: make([]*domain.Workout, 0),
	}
}

// PublishWorkoutStats mocks the PublishWorkoutStats method of AMQPPublisher
func (m *MockAMQPPublisher) PublishWorkoutStats(workoutStats *domain.Workout) error {
	// In the mock, we just store the workoutStats for verification in tests

	var challengeStatsDTO = challengeStatsDTO{
		PlayerID:        workoutStats.PlayerID,
		WorkoutEnd:      workoutStats.EndedAt,
		EnemiesFought:   workoutStats.Fights,
		EnemiesEscaped:  workoutStats.Escapes,
		DistanceCovered: workoutStats.DistanceCovered,
	}
	m.PublishedWorkouts = append(m.PublishedWorkouts, workoutStats)
	logger.Info("Workout statistics published to Challenge Manager", zap.Any("Stats", challengeStatsDTO))
	return nil // Return nil to simulate successful execution
}
