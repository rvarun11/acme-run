package httphandler

import (
	"net/http"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WorkoutHTTPHandler struct {
	gin *gin.Engine
	svc *services.WorkoutService
}

func NewWorkoutHTTPHandler(gin *gin.Engine, workoutSvc *services.WorkoutService) *WorkoutHTTPHandler {
	return &WorkoutHTTPHandler{
		gin: gin,
		svc: workoutSvc,
	}
}

func (handler *WorkoutHTTPHandler) InitRouter() {

	router := handler.gin.Group("/api/v1")

	router.POST("/workout", handler.StartWorkout)
	router.PUT("/workout/:workoutId", handler.StopWorkout)

	router.GET("/workout/:workoutId/options", handler.GetWorkoutOptions)
	router.POST("/workout/:workoutId/options", handler.StartWorkoutOption)
	router.PATCH("/workout/:workoutId/options", handler.StopWorkoutOption)

	router.GET("workout/distance", handler.GetDistance)
	router.GET("workout/shelters", handler.GetShelters)
	router.GET("workout/escapes", handler.GetEscapes)
	router.GET("workout/fights", handler.GetFights)
}

func (h *WorkoutHTTPHandler) StartWorkout(ctx *gin.Context) {

	var startWorkout StartWorkout
	if err := ctx.ShouldBindJSON(&startWorkout); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	workout, err := domain.NewWorkout(startWorkout.PlayerID, startWorkout.TrailID, startWorkout.HRMId, startWorkout.HRMConnected, startWorkout.HardCoreMode)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	linkURL, err := h.svc.Start(&workout, startWorkout.HRMId, startWorkout.HRMConnected)
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

func (h *WorkoutHTTPHandler) StopWorkout(ctx *gin.Context) {
	// Retrieve workoutId from the path parameter
	workoutId := ctx.Param("workoutId")

	workoutID, err := uuid.Parse(workoutId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Workout ID : " + workoutID.String(),
		})
		return
	}

	w, err := h.svc.Stop(workoutID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, w)
}

func (h *WorkoutHTTPHandler) GetWorkoutOptions(ctx *gin.Context) {
	// Retrieve workoutId from the path parameter
	workoutIdStr := ctx.Param("workoutId")

	workoutID, err := uuid.Parse(workoutIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Workout ID : " + workoutID.String(),
		})
		return
	}

	workoutOptions, err := h.svc.GetWorkoutOptions(workoutID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, workoutOptions)
}

func (h *WorkoutHTTPHandler) StartWorkoutOption(ctx *gin.Context) {
	// Retrieve workoutId from the path parameter
	workoutIdStr := ctx.Param("workoutId")

	workoutID, err := uuid.Parse(workoutIdStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Workout ID : " + workoutID.String(),
		})
		return
	}

	var startWorkout StartWorkoutOption
	if err := ctx.ShouldBindJSON(&startWorkout); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Bad POST Request"})
		return
	}

	err = h.svc.StartWorkoutOption(workoutID, uint8(startWorkout.Option))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "workout option started successfully"})
}

func (h *WorkoutHTTPHandler) StopWorkoutOption(ctx *gin.Context) {
	// Retrieve workoutId from the path parameter
	workoutIdStr := ctx.Param("workoutId")

	workoutID, err := uuid.Parse(workoutIdStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Workout ID : " + workoutID.String(),
		})
		return
	}

	err = h.svc.StopWorkoutOption(workoutID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Workout option stopped successfully"})
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

func (h *WorkoutHTTPHandler) GetDistance(ctx *gin.Context) {
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

func (h *WorkoutHTTPHandler) GetShelters(ctx *gin.Context) {
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

func (h *WorkoutHTTPHandler) GetEscapes(ctx *gin.Context) {
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

func (h *WorkoutHTTPHandler) GetFights(ctx *gin.Context) {
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
