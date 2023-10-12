package handler

import (
	"net/http"

	"github.com/CAS735-F23/macrun-teamvs_/workout/internal/core/services"

	"github.com/CAS735-F23/macrun-teamvs_/workout/internal/core/domain"

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

func (h *HTTPHandler) StartWorkout(ctx *gin.Context) {
	var p domain.Workout
	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	// TODO: Should this be here or in the service?
	// Keeping it here for now so that uuid can be sent back
	workout, err := domain.NewWorkout(p)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
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

// TODO: To be implemented
func (h *HTTPHandler) UpdateWorkout(ctx *gin.Context) {

}
