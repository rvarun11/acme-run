package handler

import (
	"net/http"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/dto"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/services"

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

func (h *HTTPHandler) StartWorkout(ctx *gin.Context) {

	//TODO Error Handling

	var startWorkout dto.StartWorkout
	if err := ctx.ShouldBindJSON(&startWorkout); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	workout, err := domain.NewWorkout(startWorkout.PlayerID, startWorkout.TrailID, startWorkout.HRMId, startWorkout.HRMConnected)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// TODO: The two error handling are the same, it can be refactored
	linkURL, err := h.svc.Start(&workout)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":                   workout.WorkoutID,
		"workout options link": linkURL,
		"message":              "New workout created successfully",
	})
}

func (h *HTTPHandler) StopWorkout(ctx *gin.Context) {

	var stopWorkout dto.StopWorkout
	if err := ctx.ShouldBindJSON(&stopWorkout); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid Workout ID",
		})
		return
	}

	w, err := h.svc.Stop(stopWorkout.WorkoutID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusAccepted, w)
}

func (h *HTTPHandler) GetWorkoutOptions(ctx *gin.Context) {
	workoutID, err := parseUUID(ctx, "workoutID")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Workout ID" + workoutID.String(),
		})
		return
	}

	workoutOptions, err := h.svc.GetWorkoutOptions(workoutID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, workoutOptions)
}

func (h *HTTPHandler) StartWorkoutOption(ctx *gin.Context) {
	var startWorkout dto.StartWorkoutOption

	if err := ctx.ShouldBindJSON(&startWorkout); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": "Bad POST Request",
		})
		return
	}

	err := h.svc.StartWorkoutOption(startWorkout.WorkoutID, uint8(startWorkout.Option))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "workout option started successfully",
	})
}

func (h *HTTPHandler) StopWorkoutOption(ctx *gin.Context) {

	var stopWorkoutOption dto.StopWorkoutOption

	if err := ctx.ShouldBindJSON(&stopWorkoutOption); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": "Bad POST Request",
		})
		return
	}

	err := h.svc.StopWorkoutOption(stopWorkoutOption.WorkoutID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Workout option stopped successfully",
	})
}

func parseUUID(ctx *gin.Context, paramName string) (uuid.UUID, error) {
	uuidStr := ctx.Query(paramName)
	uuidValue, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuidValue, nil
}

func parseTime(ctx *gin.Context, paramName string, layout string) (time.Time, error) {
	timeStr := ctx.Params.ByName(paramName)
	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}

func (h *HTTPHandler) GetDistance(ctx *gin.Context) {
	workoutID, err := parseUUID(ctx, "workoutID")
	var distance float64

	if err == nil {
		distance, err = h.svc.GetDistanceById(workoutID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
	} else {
		playerID, err := parseUUID(ctx, "playerID")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid player ID",
			})
			return
		}

		startDate, err := parseTime(ctx, "startDate", time.RFC3339)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid startDate",
			})
			return
		}

		endDate, err := parseTime(ctx, "endDate", time.RFC3339)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid endDate",
			})
			return
		}

		distance, err = h.svc.GetDistanceCoveredBetweenDates(playerID, startDate, endDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"workoutID":       workoutID,
		"distanceCovered": distance,
	})
}

func (h *HTTPHandler) GetShelters(ctx *gin.Context) {
	var shelterCount uint16
	workoutID, workoutIDErr := parseUUID(ctx, "workoutID")
	playerID, playerIDErr := parseUUID(ctx, "playerID")

	if workoutIDErr == nil {
		// If workoutID is provided, call GetSheltersTakenById
		var err error
		shelterCount, err = h.svc.GetSheltersTakenById(workoutID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
	} else if playerIDErr == nil {
		// If playerID, startDate, and endDate are provided, call GetSheltersTakenBetweenDates
		startDate, startDateErr := parseTime(ctx, "startDate", time.RFC3339)
		endDate, endDateErr := parseTime(ctx, "endDate", time.RFC3339)

		if startDateErr != nil || endDateErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid startDate or endDate",
			})
			return
		}

		var err error
		shelterCount, err = h.svc.GetSheltersTakenBetweenDates(playerID, startDate, endDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
	} else {
		// Handle the case where neither workoutID nor playerID is provided
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid workoutID or playerID",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"shelterCount": shelterCount,
	})
}

func (h *HTTPHandler) GetEscapes(ctx *gin.Context) {
	var escapeCount uint16
	workoutID, workoutIDErr := parseUUID(ctx, "workoutID")
	playerID, playerIDErr := parseUUID(ctx, "playerID")

	if workoutIDErr == nil {
		// If workoutID is provided, call GetEscapesMadeById
		var err error
		escapeCount, err = h.svc.GetEscapesMadeById(workoutID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
	} else if playerIDErr == nil {
		// If playerID, startDate, and endDate are provided, call GetEscapesMadeBetweenDates
		startDate, startDateErr := parseTime(ctx, "startDate", time.RFC3339)
		endDate, endDateErr := parseTime(ctx, "endDate", time.RFC3339)

		if startDateErr != nil || endDateErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid startDate or endDate",
			})
			return
		}

		var err error
		escapeCount, err = h.svc.GetEscapesMadeBetweenDates(playerID, startDate, endDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
	} else {
		// Handle the case where neither workoutID nor playerID is provided
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid workoutID or playerID",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"escapeCount": escapeCount,
	})
}

func (h *HTTPHandler) GetFights(ctx *gin.Context) {
	var fightCount uint16
	workoutID, workoutIDErr := parseUUID(ctx, "workoutID")
	playerID, playerIDErr := parseUUID(ctx, "playerID")

	if workoutIDErr == nil {
		// If workoutID is provided, call GetFightsFoughtById
		var err error
		fightCount, err = h.svc.GetFightsFoughtById(workoutID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
	} else if playerIDErr == nil {
		// If playerID, startDate, and endDate are provided, call GetFightsFoughtBetweenDates
		startDate, startDateErr := parseTime(ctx, "startDate", time.RFC3339)
		endDate, endDateErr := parseTime(ctx, "endDate", time.RFC3339)

		if startDateErr != nil || endDateErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid startDate or endDate",
			})
			return
		}

		var err error
		fightCount, err = h.svc.GetFightsFoughtBetweenDates(playerID, startDate, endDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
	} else {
		// Handle the case where neither workoutID nor playerID is provided
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid workoutID or playerID",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"fightCount": fightCount,
	})
}
