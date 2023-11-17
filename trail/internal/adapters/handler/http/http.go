package httphandler

import (
	"net/http"
	"strconv"

	"github.com/CAS735-F23/macrun-teamvsl/trail/internal/core/services"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type TrailManagerServiceHTTPHandler struct {
	gin *gin.Engine
	tvc *services.TrailManagerService
}

func NewTrailManagerServiceHTTPHandler(gin *gin.Engine, tmSvc *services.TrailManagerService) *TrailManagerServiceHTTPHandler {
	return &TrailManagerServiceHTTPHandler{
		gin: gin,
		tvc: tmSvc,
	}
}

func (handler *TrailManagerServiceHTTPHandler) InitRouter() {

	router := handler.gin.Group("/api/v1")

	// router.POST("/trail_manager", handler.ConnectToTrailManager)
	// router.PUT("/trail_manager", handler.CloseTrailManager)
	router.GET("/trail_manager/trail/", handler.GetClosestTrail)
	router.GET("/trail_manager/trail/info", handler.GetTrailLocationInfo)
	router.POST("/trail_manager/trail/create", handler.CreateTrail)
	router.PUT("/trail_manager/trail/update", handler.UpdateTrail)
	router.PUT("/trail_manager/trail/delete", handler.DeleteTrail)

	router.GET("/trail_manager/shelter/check_status", handler.CheckShelterStatus)
	router.GET("/trail_manager/shelter/info", handler.GetShelterLocationInfo)
	router.POST("/trail_manager/shelter/create", handler.CreateShelter)
	router.PUT("/trail_manager/shelter/update", handler.UpdateShelter)
	router.PUT("/trail_manager/shelter/delete", handler.DeleteShelter)

	router.POST("/trail_manager/zone/create", handler.CreateZone)
	router.PUT("/trail_manager/zone/update", handler.UpdateZone)
	router.PUT("/trail_manager/zone/delete", handler.DeleteZone)

}

func parseUUID(ctx *gin.Context, paramName string) (uuid.UUID, error) {
	uuidStr := ctx.Query(paramName)
	uuidValue, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuidValue, nil
}

func (s *TrailManagerServiceHTTPHandler) CreateTrail(ctx *gin.Context) {

	var trailDataInstance TrailDTO
	if err := ctx.ShouldBindJSON(&trailDataInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	id := trailDataInstance.TrailID
	name := trailDataInstance.TrailName
	startLongitude := trailDataInstance.StartLongitude
	startLatitude := trailDataInstance.StartLatitude
	endLongitude := trailDataInstance.EndLongitude
	endLatitude := trailDataInstance.EndLatitude
	zId := trailDataInstance.ZoneID

	_, err := s.tvc.CreateTrail(id, name, zId, startLatitude, startLongitude, endLatitude, endLongitude)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trail"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "trail created successfully"})

}

func (s *TrailManagerServiceHTTPHandler) UpdateTrail(c *gin.Context) {

	id := c.Query("id")
	name := c.Query("name")
	zId := c.Query("zone_id")
	startLongitude := c.Query("start_longitude")
	startLatitude := c.Query("start_latitude")
	endLongitude := c.Query("end_longitude")
	endLatitude := c.Query("end_latitude")

	err := s.tvc.UpdateTrail(id, name, zId, startLatitude, startLongitude, endLatitude, endLongitude, shelterId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to update trail"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"status": "success", "message": "updated trail"})

}

func (t *TrailManagerServiceHTTPHandler) DeleteTrail(ctx *gin.Context) {
	id, _ := parseUUID(ctx, "trail_id")
	err := t.tvc.DeleteTrail(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to delete trail"})
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{"status": "success", "message": "deleted trail"})
}

func (t *TrailManagerServiceHTTPHandler) GetClosestTrail(ctx *gin.Context) {

	longitudeStr := ctx.Query("longitude")
	latitudeStr := ctx.Query("latitude")
	zId, _ := parseUUID(ctx, "zone_id")

	// Convert query parameters to float64
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)

	latitude, _ := strconv.ParseFloat(latitudeStr, 64)

	// Assuming you have a method to find the closest trails by coordinates
	closestTrail, err := t.tvc.GetClosestTrail(zId, longitude, latitude)
	if err != nil || closestTrail == uuid.Nil {
		// Handle possible errors, such as no trails being found
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to retrieve closest trail"})
		return
	}

	// Respond with the ID of the closest trail
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "trail_id": closestTrail})
}

func (t *TrailManagerServiceHTTPHandler) GetTrailLocationInfo(ctx *gin.Context) {

	tId, _ := parseUUID(ctx, "trail_id")

	// Assuming you have a method to find the closest trails by coordinates
	trail, err := t.tvc.GetTrailByID(tId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to retrieve trail info"})
		return
	}

	// Respond with the ID of the closest trail
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "start_longitude": trail.StartLongitude, "start_latitude": trail.StartLatitude,
		"end_longitude": trail.EndLongitude, "end_latitude": trail.EndLatitude})
}

// shelters
func (t *TrailManagerServiceHTTPHandler) CreateShelter(ctx *gin.Context) {

	var shelterDataInstance ShelterDTO
	if err := ctx.ShouldBindJSON(&shelterDataInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	id := shelterDataInstance.ShelterID
	name := shelterDataInstance.ShelterName
	longitude := shelterDataInstance.Longitude
	latitude := shelterDataInstance.Latitude
	tId := shelterDataInstance.TrailID

	_, err := t.tvc.CreateShelter(id, name, tId, true, longitude, latitude)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to create shelter"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "shelter created successfully"})
}

func (t *TrailManagerServiceHTTPHandler) UpdateShelter(ctx *gin.Context) {

	sId, _ := parseUUID(ctx, "shelter_id")
	name := ctx.Query("name")
	longitudeStr := ctx.Query("longitude")
	latitudeStr := ctx.Query("latitude")
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)
	latitude, _ := strconv.ParseFloat(latitudeStr, 64)
	availability, _ := strconv.ParseBool("shelter_availability")
	tId, _ := parseUUID(ctx, "trail_id")

	err := t.tvc.UpdateShelter(sId, name, tId, availability, longitude, latitude)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create shelter"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "trail updated successfully"})
}

func (t *TrailManagerServiceHTTPHandler) DeleteShelter(ctx *gin.Context) {

	id, _ := parseUUID(ctx, "shelter_id")

	err := t.tvc.DeleteShelter(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to create shelter"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "shelter deleted successfully"})
}

func (h *TrailManagerServiceHTTPHandler) CheckShelterStatus(ctx *gin.Context) {
	sId, _ := parseUUID(ctx, "shelter_id")
	// first check if the workout is bind to any trail manager, if not raise error
	shelterInstance, err := h.tvc.GetShelterByID(sId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "cant find shelter"})
		return
	}
	if shelterInstance.ShelterAvailability {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "available"})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "unavailable"})
		return
	}

}

func (t *TrailManagerServiceHTTPHandler) GetShelterLocationInfo(ctx *gin.Context) {

	id, _ := parseUUID(ctx, "shelter_id")

	// Assuming you have a method to find the closest trails by coordinates
	shelter, err := t.tvc.GetShelterByID(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "failed to retrieve shelter info"})
		return
	}

	// Respond with the ID of the closest trail
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "longitude": shelter.Longitude, "latitude": shelter.Latitude})
}

func (h *TrailManagerServiceHTTPHandler) CreateZone(ctx *gin.Context) {
	var zoneDataInstance ZoneDTO
	if err := ctx.ShouldBindJSON(&zoneDataInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	zId, err := h.tvc.CreateZone(zoneDataInstance.ZoneName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
	} else {
		zIdString := zId.String()
		ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "zone created with id", "zone_id": zIdString})
	}
	return
}

// 	// check if this trail has any shelter, if not return error
// 	availability, err := h.tvc.CheckTrailShelter(tId)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}
// 	var sv ShelterAvailable
// 	sv.ShelterAvailable = availability
// 	sv.WorkoutID = wId
// 	var outputData []byte
// 	var err2 error
// 	if !availability {
// 		// if there is no shelter attached to the trail, return
// 		sv.DistanceToShelter = math.MaxFloat64
// 		sv.ShelterAvailable = false
// 		sv.WorkoutID = wId
// 		outputData, _ = json.Marshal(sv)
// 		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": string(outputData)})
// 		return
// 	} else {
// 		long := ctx.Query("longitude")
// 		lat := ctx.Query("latitude")
// 		inputLongitude, _ := strconv.ParseFloat(long, 64)
// 		inputLatitude, _ := strconv.ParseFloat(lat, 64)
// 		// get current location for calculating the distance, if the user is not connected to any instance
// 		// it will raise error and return
// 		if returnedTM.CurrentLongitude == math.MaxFloat64 || returnedTM.CurrentLatitude == math.MaxFloat64 {
// 			fmt.Println(err)
// 			ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "failed to get shelter info"})
// 			return
// 		}
// 		sv.DistanceToShelter, _ = h.tvc.CalculateDistance(inputLongitude, inputLatitude, returnedTM.CurrentLongitude, returnedTM.CurrentLatitude)
// 	}
// 	outputData, err2 = json.Marshal(sv)
// 	if err2 != nil {
// 		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "failed to get shelter info"})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": string(outputData)})

// }

// func (h *TrailManagerServiceHTTPHandler) GetDistance(ctx *gin.Context) {
// 	workoutID, err := parseUUID(ctx, "workoutID")
// 	var distance float64

// 	if err == nil {
// 		distance, err = h.tvc.GetDistance()
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{
// 				"error": err,
// 			})
// 			return
// 		}
// 	} else {
// 		playerID, err := parseUUID(ctx, "playerID")
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{
// 				"error": "Invalid player ID",
// 			})
// 			return
// 		}

// 		startDate, err := parseTime(ctx, "startDate", time.RFC3339)
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{
// 				"error": "Invalid startDate",
// 			})
// 			return
// 		}

// 		endDate, err := parseTime(ctx, "endDate", time.RFC3339)
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{
// 				"error": "Invalid endDate",
// 			})
// 			return
// 		}

// 		distance, err = h.tvc.GetDistanceCoveredBetweenDates(playerID, startDate, endDate)
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{
// 				"error": err,
// 			})
// 			return
// 		}
// 	}

// 	ctx.JSON(http.StatusCreated, gin.H{
// 		"workoutID":           workoutID,
// 		"distance_to_shelter": distance,
// 	})
// }

// 	c.JSON(http.StatusCreated, gin.H{"message": "Trail created successfully"})
// }

// func (s *TrailManagerServiceHTTPHandler) UpdateTrail(c *gin.Context) {
// 	// Parse the trail ID from the URL parameter or query, depending on your routing setup
// 	trailID, err := uuid.Parse(c.Query("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trail ID format"})
// 		return
// 	}

// 	name := c.Query("name")
// 	startLongitude, err := strconv.ParseFloat(c.Query("start_longitude"), 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start longitude format"})
// 		return
// 	}
// 	startLatitude, err := strconv.ParseFloat(c.Query("start_latitude"), 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start latitude format"})
// 		return
// 	}
// 	endLongitude, err := strconv.ParseFloat(c.Query("end_longitude"), 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end longitude format"})
// 		return
// 	}
// 	endLatitude, err := strconv.ParseFloat(c.Query("end_latitude"), 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end latitude format"})
// 		return
// 	}
// 	shelterID, err := uuid.Parse(c.Query("shelter_id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shelter ID format"})
// 		return
// 	}

// 	// Call the UpdateTrailByID method from TrailRepository
// 	err = s.tvc.repoT.UpdateTrailByID(trailID, name, startLatitude, startLongitude, endLatitude, endLongitude, shelterID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update trail"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Trail updated successfully"})
// }

// func (s *TrailManagerServiceHTTPHandler) UpdateShelter(c *gin.Context) {
// 	// Parse the trail ID from the URL parameter or query, depending on your routing setup
// 	shelterID, err := uuid.Parse(c.Query("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trail ID format"})
// 		return
// 	}

// 	name := c.Query("name")
// 	Longitude, err := strconv.ParseFloat(c.Query("longitude"), 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid  longitude format"})
// 		return
// 	}
// 	Latitude, err := strconv.ParseFloat(c.Query("latitude"), 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude format"})
// 		return
// 	}

// 	// Call the UpdateTrailByID method from TrailRepository
// 	err = s.tvc.repoS.UpdateShelterByID(shelterID, name, Latitude, Longitude)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update shelter"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Shelter updated successfully"})
// }

// func (s *TrailManagerServiceHTTPHandler) GetClosestShelter(c *gin.Context) {
// 	longitudeStr := c.Query("longitude")
// 	latitudeStr := c.Query("latitude")

// 	// Convert query parameters to float64
// 	longitude, err := strconv.ParseFloat(longitudeStr, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude format"})
// 		return
// 	}

// 	latitude, err := strconv.ParseFloat(latitudeStr, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude format"})
// 		return
// 	}

// 	// Assuming you have a method to find the closest trails by coordinates
// 	closestShelterId, err := s.tvc.getShelter(longitude, latitude)
// 	if err != nil {
// 		// Handle possible errors, such as no trails being found
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve closest trail"})
// 		return
// 	}

// 	// Respond with the ID of the closest trail
// 	c.JSON(http.StatusOK, gin.H{"closestShelterID": closestShelterId})
// }
