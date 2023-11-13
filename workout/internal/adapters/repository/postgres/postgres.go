package postgres

import (
	"errors"
	"fmt"
	"time"

	logger "github.com/CAS735-F23/macrun-teamvsl/challenge_manager/log"
	"github.com/CAS735-F23/macrun-teamvsl/workout/config"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/domain"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(cfg *config.Postgres) *Repository {

	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable client_encoding=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DB_Name,
		cfg.Password,
		cfg.Encoding,
	)

	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&postgresWorkout{}, &postgresWorkoutOptions{})

	return &Repository{
		db: db,
	}
}

// Repository Types

type postgresWorkout struct {
	// ID is the identifier of the Entity, the ID is shared for all sub domains
	WorkoutID uuid.UUID `gorm:"type:uuid;primaryKey"`
	// trailId is the id of the trail player is on
	TrailID uuid.UUID `gorm:"type:uuid;unique not null"`
	// PlayerID of the player starting the workout session
	PlayerID uuid.UUID `gorm:"type:uuid;unique not null"`
	// InProgress tells whether the workout is in progress
	IsCompleted bool
	// CreatedAt is the time when the workout was started
	CreatedAt time.Time
	// Duration of the workout
	EndedAt time.Time
	// EndedAt is the time when the workout was ended
	DistanceCovered float64
	// Player Profile can be either 'cardio' or 'strength'
	Profile string
	// HardcoreMode is the difficulty level chosen by the player
	HardcoreMode bool
	// Shelters taken for a given workout
	Shelters uint8
	// Fights fought in a given workout
	Fights uint8
	// Escapes made in a given workout
	Escapes uint8
}

type postgresWorkoutOptions struct {
	// ID is the identifier of the Entity, the ID is shared for all sub domains
	WorkoutID uuid.UUID `gorm:"type:uuid;primary_key"`
	// WorkoutOptions a uint8 integer, uses bit representation for available workout options
	CurrentWorkoutOption int8
	// FightsPushDown Ranking bool
	FightsPushDown bool
	// Is WorkoutOption Active
	IsWorkoutOptionActive bool
	// Distance to Shelter
	DistanceToShelter float64
}

func toWorkoutAggregate(pworkout *postgresWorkout) *domain.Workout {

	return &domain.Workout{
		WorkoutID:       pworkout.WorkoutID,
		TrailID:         pworkout.TrailID,
		PlayerID:        pworkout.PlayerID,
		IsCompleted:     pworkout.IsCompleted,
		CreatedAt:       pworkout.CreatedAt,
		EndedAt:         pworkout.EndedAt,
		DistanceCovered: pworkout.DistanceCovered,
		Profile:         pworkout.Profile,
		HardcoreMode:    pworkout.HardcoreMode,
		Shelters:        pworkout.Shelters,
		Fights:          pworkout.Fights,
		Escapes:         pworkout.Escapes,
	}
}

func toWorkoutPostgres(workout *domain.Workout) *postgresWorkout {

	return &postgresWorkout{
		WorkoutID:       workout.WorkoutID,
		TrailID:         workout.TrailID,
		PlayerID:        workout.PlayerID,
		IsCompleted:     workout.IsCompleted,
		CreatedAt:       workout.CreatedAt,
		EndedAt:         workout.EndedAt,
		DistanceCovered: workout.DistanceCovered,
		Profile:         workout.Profile,
		HardcoreMode:    workout.HardcoreMode,
		Shelters:        workout.Shelters,
		Fights:          workout.Fights,
		Escapes:         workout.Escapes,
	}
}

func toWorkoutOptionsAggregate(pworkoutOptions *postgresWorkoutOptions) *domain.WorkoutOptions {

	return &domain.WorkoutOptions{
		WorkoutID:             pworkoutOptions.WorkoutID,
		CurrentWorkoutOption:  pworkoutOptions.CurrentWorkoutOption,
		FightsPushDown:        pworkoutOptions.FightsPushDown,
		IsWorkoutOptionActive: pworkoutOptions.IsWorkoutOptionActive,
	}
}

func toWorkoutOptionsPostgres(workoutOptions *domain.WorkoutOptions) *postgresWorkoutOptions {

	return &postgresWorkoutOptions{
		WorkoutID:             workoutOptions.WorkoutID,
		CurrentWorkoutOption:  workoutOptions.CurrentWorkoutOption,
		FightsPushDown:        workoutOptions.FightsPushDown,
		IsWorkoutOptionActive: workoutOptions.IsWorkoutOptionActive,
		DistanceToShelter:     workoutOptions.DistanceToShelter,
	}
}

// Repository Functions

func (r *Repository) Create(workout *domain.Workout, workoutOptions *domain.WorkoutOptions) error {
	logger.Info("DEBUG-----CREATE ROW IN REPO")

	pworkout := &postgresWorkout{
		WorkoutID:       workout.WorkoutID,
		TrailID:         workout.TrailID,
		PlayerID:        workout.PlayerID,
		IsCompleted:     workout.IsCompleted,
		CreatedAt:       workout.CreatedAt,
		EndedAt:         workout.EndedAt,
		DistanceCovered: workout.DistanceCovered,
		Profile:         workout.Profile,
		HardcoreMode:    workout.HardcoreMode,
		Shelters:        workout.Shelters,
		Fights:          workout.Fights,
		Escapes:         workout.Escapes,
	}

	pworkoutOptions := &postgresWorkoutOptions{
		WorkoutID:             workoutOptions.WorkoutID,
		IsWorkoutOptionActive: workoutOptions.IsWorkoutOptionActive,
		CurrentWorkoutOption:  workoutOptions.CurrentWorkoutOption,
		FightsPushDown:        workoutOptions.FightsPushDown,
		DistanceToShelter:     workoutOptions.DistanceToShelter,
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		res := tx.First(&pworkout, "workout_id = ?", workout.WorkoutID)
		if res.Error != nil {
			// If error occurs in finding the workout, return the error.
			return res.Error
		}

		if res.RowsAffected == 0 {
			if err := tx.Create(&pworkout).Error; err != nil {
				// Log and return error if the creation fails.
				logger.Debug("FAILED TO CREATE WORKOUT", zap.String("error", err.Error()))
				return err
			}
		}

		// Create 'pworkoutOptions' if 'pworkout' creation was successful or it already existed.
		if err := tx.Create(&pworkoutOptions).Error; err != nil {
			// Log and return error if the creation fails.
			logger.Debug("FAILED TO CREATE WORKOUT OPTIONS", zap.String("error", err.Error()))
			return err
		}
		return nil
	})

	// TODO Define Errors in Repo Interface file and return them instead of this
	if err != nil {
		return err
	}

	return err
}

func (r *Repository) GetWorkout(workoutID uuid.UUID) (*domain.Workout, error) {
	var pworkout postgresWorkout

	res := r.db.First(&pworkout, "workout_id = ?", workoutID)

	if res.Error != nil {
		return &domain.Workout{}, res.Error
	}

	workout := toWorkoutAggregate(&pworkout)
	return workout, nil
}

func (r *Repository) GetWorkoutOptions(workoutID uuid.UUID) (*domain.WorkoutOptions, error) {
	var pworkoutOptions postgresWorkoutOptions

	res := r.db.First(&pworkoutOptions, "workout_id = ?", workoutID)

	if res.Error != nil {
		return &domain.WorkoutOptions{}, res.Error
	}

	workoutOptions := toWorkoutOptionsAggregate(&pworkoutOptions)
	return workoutOptions, nil
}

func (r *Repository) UpdateWorkout(workout *domain.Workout) (*domain.Workout, error) {

	pworkout := toWorkoutPostgres(workout)

	if err := r.db.Save(&pworkout).Error; err != nil {
		return &domain.Workout{}, err
	}

	return toWorkoutAggregate(pworkout), nil
}

func (r *Repository) UpdateWorkoutOptions(workoutOptions *domain.WorkoutOptions) (*domain.WorkoutOptions, error) {

	pworkoutOptions := toWorkoutOptionsPostgres(workoutOptions)

	if err := r.db.Save(&pworkoutOptions).Error; err != nil {
		return &domain.WorkoutOptions{}, err
	}

	return toWorkoutOptionsAggregate(pworkoutOptions), nil
}

func (r *Repository) DeleteWorkoutOptions(workoutID uuid.UUID) error {
	// Define a variable to store the workout options record
	var workoutOptions postgresWorkoutOptions

	// Use GORM to find the record with the specified WorkoutID
	if err := r.db.Where("workout_id = ?", workoutID).First(&workoutOptions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// The record does not exist, no need to delete
			return nil
		}
		// Return any other error that occurred during the query
		return err
	}

	// The record exists, so use GORM to delete it
	if err := r.db.Delete(&workoutOptions).Error; err != nil {
		// Return any error that occurred during the deletion
		return err
	}

	// No errors occurred, return nil to indicate success
	return nil
}

func (r *Repository) GetDistanceByID(workoutID uuid.UUID) (float64, error) {
	var distanceCovered = 0.0

	err := r.db.Table("postgres_workouts").
		Where("workout_id = ?", workoutID).
		Pluck("distance_covered", &distanceCovered).
		Error

	// TODO Define Errors in Repo Interface file and return them instead of this
	if err != nil {
		return distanceCovered, err
	}

	return distanceCovered, err
}

func (r *Repository) GetDistanceCoveredBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (float64, error) {
	totalDistanceCovered := 0.0

	err := r.db.Table("postgres_workouts").
		Where("player_id = ? AND created_at >= ? AND ended_at <= ?", playerID, startDate, endDate).
		Select("COALESCE(SUM(DistanceCovered), 0)").Row().
		Scan(&totalDistanceCovered)

	// TODO Define Errors in Repo Interface file and return them instead of this
	if err != nil {
		return totalDistanceCovered, err
	}

	return totalDistanceCovered, err
}

func (r *Repository) GetEscapesMadeByID(workoutID uuid.UUID) (uint16, error) {
	escapesMade := 0

	err := r.db.Table("postgres_workouts").
		Where("workout_id = ?", workoutID).
		Pluck("escapes", &escapesMade).
		Error

	if err != nil {
		return uint16(escapesMade), err
	}

	return uint16(escapesMade), nil
}

func (r *Repository) GetEscapesMadeBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (uint16, error) {
	totalEscapesMade := 0

	err := r.db.Table("postgres_workouts").
		Where("player_id = ? AND created_at >= ? AND ended_at <= ?", playerID, startDate, endDate).
		Select("COALESCE(SUM(escapes), 0)").Row().
		Scan(&totalEscapesMade)

	if err != nil {
		return uint16(totalEscapesMade), err
	}

	return uint16(totalEscapesMade), nil
}

func (r *Repository) GetFightsFoughtByID(workoutID uuid.UUID) (uint16, error) {
	fightsFought := 0

	err := r.db.Table("postgres_workouts").
		Where("workout_id = ?", workoutID).
		Pluck("fights", &fightsFought).
		Error

	if err != nil {
		return uint16(fightsFought), err
	}

	return uint16(fightsFought), nil
}

func (r *Repository) GetFightsFoughtBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (uint16, error) {
	totalFightsFought := 0

	err := r.db.Table("postgres_workouts").
		Select("COALESCE(SUM(fights), 0)").Row().
		Scan(&totalFightsFought)

	if err != nil {
		return uint16(totalFightsFought), err
	}

	return uint16(totalFightsFought), nil
}

func (r *Repository) GetSheltersTakenByID(workoutID uuid.UUID) (uint16, error) {
	sheltersTaken := 0

	err := r.db.Table("postgres_workouts").
		Where("workout_id = ?", workoutID).
		Pluck("shelters", &sheltersTaken).
		Error

	if err != nil {
		return uint16(sheltersTaken), err
	}

	return uint16(sheltersTaken), nil
}

func (r *Repository) GetSheltersTakenBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (uint16, error) {
	totalSheltersTaken := 0

	err := r.db.Table("postgres_workouts").
		Where("player_id = ? AND created_at >= ? AND ended_at <= ?", playerID, startDate, endDate).
		Select("COALESCE(SUM(shelters), 0)").Row().
		Scan(&totalSheltersTaken)

	if err != nil {
		return uint16(totalSheltersTaken), err
	}

	return uint16(totalSheltersTaken), nil
}
