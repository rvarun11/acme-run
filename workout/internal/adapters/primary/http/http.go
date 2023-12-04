package httphandler

import (
	"net/http"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WorkoutHanlder struct {
	gin *gin.Engine
	svc *services.WorkoutService
}

func NewWorkoutHanlder(gin *gin.Engine, workoutSvc *services.WorkoutService) *WorkoutHanlder {
	return &WorkoutHanlder{
		gin: gin,
		svc: workoutSvc,
	}
}

func (handler *WorkoutHanlder) InitRouter() {

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

// StartWorkout starts a new workout session for a player.
//
//	@Summary		Start a new workout session
//	@Description	This endpoint starts a new workout session for a player with the given details.
//	@Tags			workout
//	@ID				start-workout
//	@Accept			json
//	@Produce		json
//	@Param			workout	body	StartWorkout	true	"Details of the workout to start"
//	@Success		201		"Successfully started workout session"
//	@Failure		400		"Bad Request with error details"
//	@Router			/api/v1/workout [post]
func (h *WorkoutHanlder) StartWorkout(ctx *gin.Context) {

	var startWorkout StartWorkout
	if err := ctx.ShouldBindJSON(&startWorkout); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request parameters",
		})
		return
	}

	workout, err := domain.NewWorkout(startWorkout.PlayerID, startWorkout.TrailID, startWorkout.HRMId, startWorkout.HRMConnected, startWorkout.HardCoreMode)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "new workout could not be created",
		})
		return
	}

	linkURL, err := h.svc.Start(&workout, startWorkout.HRMId, startWorkout.HRMConnected)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":    "new workout created successfully",
		"workout_id": workout.WorkoutID,
		"links": gin.H{
			"workout_options": linkURL,
		},
	})
}

// StopWorkout stops an ongoing workout session for a player.
//
//	@Summary		Stop an ongoing workout session
//	@Description	This endpoint stops the workout session for a player based on the provided workout ID.
//	@Tags			workout
//	@ID				stop-workout
//	@Accept			json
//	@Produce		json
//	@Param			workoutId	path	string	true	"ID of the workout session to stop"
//	@Success		202			"Successfully stopped workout session"
//	@Failure		400			"Bad Request with error details"
//	@Router			/api/v1/workout/{workoutId} [put]
func (h *WorkoutHanlder) StopWorkout(ctx *gin.Context) {
	// Retrieve workoutId from the path parameter
	workoutId := ctx.Param("workoutId")

	workoutID, err := uuid.Parse(workoutId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid workout id",
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

// GetWorkoutOptions retrieves available options for a workout session.
//
//	@Summary		Get workout session options
//	@Description	This endpoint retrieves the available options for a workout session based on the workout ID.
//	@Tags			workout
//	@ID				get-workout-options
//	@Accept			json
//	@Produce		json
//	@Param			workoutId	path	string	true	"ID of the workout session"
//	@Success		200			"Successfully retrieved workout options"
//	@Failure		400			"Bad Request with error details"
//	@Router			/api/v1/workout/{workoutId}/options [get]
func (h *WorkoutHanlder) GetWorkoutOptions(ctx *gin.Context) {
	// Retrieve workoutId from the path parameter
	workoutIdStr := ctx.Param("workoutId")

	workoutID, err := uuid.Parse(workoutIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid workout id",
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

// StartWorkoutOption starts a specific option for an ongoing workout session.
//
//	@Summary		Start a workout option
//	@Description	This endpoint starts a specific option for an ongoing workout session based on the workout ID and option details.
//	@Tags			workout
//	@ID				start-workout-option
//	@Accept			json
//	@Produce		json
//	@Param			workoutId	path	string				true	"ID of the workout session"
//	@Param			option		body	StartWorkoutOption	true	"Details of the workout option to start"
//	@Success		200			"Successfully started workout option"
//	@Failure		400			"Bad Request with error details"
//	@Router			/api/v1/workout/{workoutId}/options [post]
func (h *WorkoutHanlder) StartWorkoutOption(ctx *gin.Context) {
	// Retrieve workoutId from the path parameter
	workoutIdStr := ctx.Param("workoutId")

	workoutID, err := uuid.Parse(workoutIdStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid workout id",
		})
		return
	}

	var startWorkout StartWorkoutOption
	if err := ctx.ShouldBindJSON(&startWorkout); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request parameters"})
		return
	}

	option, err := h.svc.StartWorkoutOption(workoutID, startWorkout.Option)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "workout option started successfully",
		"option":  option,
	})
}

// StopWorkoutOption stops a specific option of an ongoing workout session.
//
//	@Summary		Stop a workout option
//	@Description	This endpoint stops a specific option of an ongoing workout session based on the workout ID.
//	@Tags			workout
//	@ID				stop-workout-option
//	@Accept			json
//	@Produce		json
//	@Param			workoutId	path	string	true	"ID of the workout session"
//	@Success		200			"Successfully stopped workout option"
//	@Failure		400			"Bad Request with error details"
//	@Router			/api/v1/workout/{workoutId}/options [patch]
func (h *WorkoutHanlder) StopWorkoutOption(ctx *gin.Context) {
	// Retrieve workoutId from the path parameter
	workoutIdStr := ctx.Param("workoutId")

	workoutID, err := uuid.Parse(workoutIdStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid workout id",
		})
		return
	}

	option, err := h.svc.StopWorkoutOption(workoutID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "workout option stopped successfully",
		"option":  option,
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

// GetDistance retrieves the distance covered in a workout session.
//
//	@Summary		Get distance covered in a workout
//	@Description	This endpoint retrieves the distance covered in a workout session either by workout ID or by player ID within a date range.
//	@Tags			workout
//	@ID				get-distance
//	@Accept			json
//	@Produce		json
//	@Param			workoutID	query	string	false	"ID of the workout session"
//	@Param			playerID	query	string	false	"ID of the player"
//	@Param			startDate	query	string	false	"Start date for the range (RFC3339 format)"
//	@Param			endDate		query	string	false	"End date for the range (RFC3339 format)"
//	@Success		201			"Successfully retrieved distance"
//	@Failure		400			"Bad Request with error details"
//	@Router			/api/v1/workout/distance [get]
func (h *WorkoutHanlder) GetDistance(ctx *gin.Context) {
	workoutID, err := parseUUID(ctx, "workoutID")
	var distance float64

	if err == nil {
		distance, err = h.svc.GetDistanceById(workoutID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
		playerID, err := parseUUID(ctx, "playerID")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid player id",
			})
			return
		}

		startDate, err := parseTime(ctx, "startDate", time.RFC3339)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid start date",
			})
			return
		}

		endDate, err := parseTime(ctx, "endDate", time.RFC3339)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid end date",
			})
			return
		}

		distance, err = h.svc.GetDistanceCoveredBetweenDates(playerID, startDate, endDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"workout_id":       workoutID,
		"distance_covered": distance,
	})
}

// GetShelters retrieves the number of shelters taken in a workout.
//
//	@Summary		Get shelters taken in a workout
//	@Description	This endpoint retrieves the number of shelters taken either by workout ID or between dates for a player.
//	@Tags			workout
//	@ID				get-shelters
//	@Accept			json
//	@Produce		json
//	@Param			workoutID	query	string	false	"Workout ID to fetch shelters"
//	@Param			playerID	query	string	false	"Player ID to fetch shelters between dates"
//	@Param			startDate	query	string	false	"Start date for fetching shelters"
//	@Param			endDate		query	string	false	"End date for fetching shelters"
//	@Success		201			"Successfully retrieved shelter count"
//	@Failure		400			"Bad Request with error details"
//	@Router			/api/v1/workout/shelters [get]
func (h *WorkoutHanlder) GetShelters(ctx *gin.Context) {
	var shelterCount uint16
	workoutID, workoutIDErr := parseUUID(ctx, "workoutID")
	playerID, playerIDErr := parseUUID(ctx, "playerID")

	if workoutIDErr == nil {
		// If workoutID is provided, call GetSheltersTakenById
		var err error
		shelterCount, err = h.svc.GetSheltersTakenById(workoutID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else if playerIDErr == nil {
		// If playerID, startDate, and endDate are provided, call GetSheltersTakenBetweenDates
		startDate, startDateErr := parseTime(ctx, "startDate", time.RFC3339)
		endDate, endDateErr := parseTime(ctx, "endDate", time.RFC3339)

		if startDateErr != nil || endDateErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid start date or end date",
			})
			return
		}

		var err error
		shelterCount, err = h.svc.GetSheltersTakenBetweenDates(playerID, startDate, endDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
		// Handle the case where neither workoutID nor playerID is provided
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid workout id or player id",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"shelter_count": shelterCount,
	})
}

// GetEscapes retrieves the number of escapes made in a workout.
//
//	@Summary		Get escapes made in a workout
//	@Description	This endpoint retrieves the number of escapes made either by workout ID or between dates for a player.
//	@Tags			workout
//	@ID				get-escapes
//	@Accept			json
//	@Produce		json
//	@Param			workoutID	query	string	false	"Workout ID to fetch escapes"
//	@Param			playerID	query	string	false	"Player ID to fetch escapes between dates"
//	@Param			startDate	query	string	false	"Start date for fetching escapes"
//	@Param			endDate		query	string	false	"End date for fetching escapes"
//	@Success		201			"Successfully retrieved escape count"
//	@Failure		400			"Bad Request with error details"
//	@Router			/api/v1/workout/escapes [get]
func (h *WorkoutHanlder) GetEscapes(ctx *gin.Context) {
	var escapeCount uint16
	workoutID, workoutIDErr := parseUUID(ctx, "workoutID")
	playerID, playerIDErr := parseUUID(ctx, "playerID")

	if workoutIDErr == nil {
		// If workoutID is provided, call GetEscapesMadeById
		var err error
		escapeCount, err = h.svc.GetEscapesMadeById(workoutID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else if playerIDErr == nil {
		// If playerID, startDate, and endDate are provided, call GetEscapesMadeBetweenDates
		startDate, startDateErr := parseTime(ctx, "startDate", time.RFC3339)
		endDate, endDateErr := parseTime(ctx, "endDate", time.RFC3339)

		if startDateErr != nil || endDateErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid start date or end date",
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
			"error": "invalid workout id or player id",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"escape_count": escapeCount,
	})
}

// GetFights retrieves the number of fights fought in a workout.
//
//	@Summary		Get fights fought in a workout
//	@Description	This endpoint retrieves the number of fights fought either by workout ID or between dates for a player.
//	@Tags			workout
//	@ID				get-fights
//	@Accept			json
//	@Produce		json
//	@Param			workoutID	query	string	false	"Workout ID to fetch fights"
//	@Param			playerID	query	string	false	"Player ID to fetch fights between dates"
//	@Param			startDate	query	string	false	"Start date for fetching fights"
//	@Param			endDate		query	string	false	"End date for fetching fights"
//	@Success		201			"Successfully retrieved fight count"
//	@Failure		400			"Bad Request with error details"
//	@Router			/api/v1/workout/fights [get]
func (h *WorkoutHanlder) GetFights(ctx *gin.Context) {
	var fightCount uint16
	workoutID, workoutIDErr := parseUUID(ctx, "workoutID")
	playerID, playerIDErr := parseUUID(ctx, "playerID")

	if workoutIDErr == nil {
		// If workoutID is provided, call GetFightsFoughtById
		var err error
		fightCount, err = h.svc.GetFightsFoughtById(workoutID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else if playerIDErr == nil {
		// If playerID, startDate, and endDate are provided, call GetFightsFoughtBetweenDates
		startDate, startDateErr := parseTime(ctx, "startDate", time.RFC3339)
		endDate, endDateErr := parseTime(ctx, "endDate", time.RFC3339)

		if startDateErr != nil || endDateErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid start date or end date",
			})
			return
		}

		var err error
		fightCount, err = h.svc.GetFightsFoughtBetweenDates(playerID, startDate, endDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
		// Handle the case where neither workoutID nor playerID is provided
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid workout id or player id",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"fight_count": fightCount,
	})
}

// DeleteWorkout deletes a specified workout session.
//
//	@Summary		Delete a workout session
//	@Description	This endpoint deletes a workout session based on the provided workout ID.
//	@Tags			workout
//	@ID				delete-workout
//	@Accept			json
//	@Produce		json
//	@Param			workoutId	path	string	true	"ID of the workout session to delete"
//	@Success		200			"Successfully deleted workout session"
//	@Failure		400			"Bad Request with error details"
//	@Failure		404			"Workout session not found"
//	@Router			/api/v1/workout/{workoutId} [delete]
func (h *WorkoutHanlder) DeleteWorkout(ctx *gin.Context) {
	workoutId := ctx.Param("workoutId")

	// Parse the UUID from the workoutId, handle error if invalid
	_, err := uuid.Parse(workoutId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid workout id",
		})
		return
	}

	// TO BE IMPLEMENTED
	// Call the service layer to delete the workout
	/*err = h.svc.DeleteWorkout(workoutID)
	if err != nil {
		// Assuming the error is because the workout was not found
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Workout session not found",
		})
		return
	}*/

	ctx.JSON(http.StatusOK, gin.H{
		"message": "workout session deleted successfully",
	})
}
