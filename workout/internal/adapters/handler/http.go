package handler

import (
	"net/http"

	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/services"

	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HTTPHandler struct {
	svc services.WorkoutService
}

func NewHTTPHandler(WorkoutService services.WorkoutService) *HTTPHandler {
	return &HTTPHandler{
		svc: WorkoutService,
	}
}

func (h *HTTPHandler) ListWorkouts(ctx *gin.Context) {

	workouts, err := h.svc.List()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, workouts)

}

// TODO: Move all the logic to Start() of WorkoutService
func (h *HTTPHandler) StartWorkout(ctx *gin.Context) {

	// TODO: Temp DTO
	type tempDTO struct {
		TrailID uuid.UUID `json:"trailID"`
		// PlayerID of the player starting the workout session
		PlayerID uuid.UUID `json:"playerID"`
		// TODO: Remove this field. If hrmId == nil then it's the same thing
		HRMConnected bool `json:"hrmConnected"`

		HRMId uuid.UUID `json:"hrmID"`
	}
	var w domain.Workout
	var tempDTOInstance tempDTO
	if err := ctx.ShouldBindJSON(&tempDTOInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	w.PlayerID = tempDTOInstance.PlayerID
	w.TrailID = tempDTOInstance.TrailID
	// TODO: Should this be here or in the service?
	// @Samkith : I think in the service
	// Keeping it here for now so that uuid can be sent back
	// TODO: Only one workout can be active for a given player at any given time, handle that case
	workout, err := domain.NewWorkout(w)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// Send request to tie HRM to Workout
	if tempDTOInstance.HRMConnected {
		// TODO: Move to the service later along with the above code
		StartHRM(tempDTOInstance.HRMId, workout.ID)
	}

	// TODO: The two error handling are the same, it can be refactored
	err = h.svc.Start(workout)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":      workout.ID,
		"message": "New workout created successfully",
	})
}

func (h *HTTPHandler) StopWorkout(ctx *gin.Context) {

	// TODO: Temp DTO
	type tempDTO struct {
		// WorkoutID
		WorkoutID uuid.UUID `json:"workoutID"`
	}

	var tempDTOInstance tempDTO
	if err := ctx.ShouldBindJSON(&tempDTOInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	var w *domain.Workout
	var err error
	// TODO: The two error handling are the same, it can be refactored
	w, err = h.svc.Stop(tempDTOInstance.WorkoutID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusAccepted, w)
}

func (h *HTTPHandler) GetWorkout(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	workout, err := h.svc.Get(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, workout)
}

// TODO: This should probably be in the rabbitmq handler
func (h *HTTPHandler) UpdateWorkout(ctx *gin.Context) {

}
