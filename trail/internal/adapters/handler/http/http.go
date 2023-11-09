package httphandler

import (
	"net/http"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/trail/internal/core/services"

	"github.com/gin-gonic/gin"
)

type TrailManagerServiceHTTPHandler struct {
	gin *gin.Engine
	tvc *services.TrailManagerService
}

func NewTrailManagerServiceHTTPHandler(gin *gin.Engine, tmSvc *services.TrailManagerService) *WorkoutHTTPHandler {
	return &TrailManagerServiceHTTPHandler{
		gin: gin,
		tvc: tmSvc,
	}
}

func (handler *TrailManagerServiceHTTPHandler) InitRouter() {

	router := handler.gin.Group("/api/v1")

	router.POST("/trail", handler.CreateTrail)
	router.PUT("/trail", handler.UpdateTrail)
	router.GET("/trail", handler.GetClosestTrailHandler)

	router.POST("/shelter", handler.CreateShelter)
	router.PUT("/shelter", handler.UpdateShelter)
	router.GET("/shelter", handler.GetClosestShelterHandler)

}

func parseUUID(ctx *gin.Context, paramName string) (uuid.UUID, error) {
	uuidStr := ctx.Query(paramName)
	uuidValue, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuidValue, nil
}

func (h *NewTrailManagerServiceHTTPHandler) GetDistance(ctx *gin.Context) {
	workoutID, err := parseUUID(ctx, "workoutID")
	var distance float64

	if err == nil {
		distance, err = h.svc.GetDistance()
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
		"workoutID":           workoutID,
		"distance_to_shelter": distance,
	})
}


func (s *NewTrailManagerServiceHTTPHandler) CreateTrail(c *gin.Context) {

	tid : = c.Query("id")
	name := c.Query("name")
	startLongitude := c.Query("start_longitude")
	startLatitude := c.Query("start_latitude")
	endLongitude := c.Query("end_longitude")
	endLatitude := c.Query("end_latitude")
	shelterId := c.Query("shelter_id")
	// Now, create the Trail entity using the parsed parameters
	// This is a placeholder; you'll need to implement the actual creation logic
	err = s.tvc.repoT.CreateTrail(tid, name, startLatitude, startLongitude, endLatitude, endLongitude, shelterId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trail"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Trail created successfully"})
}

func (s *TrailManagerServiceHTTPHandler) UpdateTrail(c *gin.Context) {
	// Parse the trail ID from the URL parameter or query, depending on your routing setup
	trailID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trail ID format"})
		return
	}

	name := c.Query("name")
	startLongitude, err := strconv.ParseFloat(c.Query("start_longitude"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start longitude format"})
		return
	}
	startLatitude, err := strconv.ParseFloat(c.Query("start_latitude"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start latitude format"})
		return
	}
	endLongitude, err := strconv.ParseFloat(c.Query("end_longitude"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end longitude format"})
		return
	}
	endLatitude, err := strconv.ParseFloat(c.Query("end_latitude"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end latitude format"})
		return
	}
	shelterID, err := uuid.Parse(c.Query("shelter_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shelter ID format"})
		return
	}

	// Call the UpdateTrailByID method from TrailRepository
	err = s.tvc.repoT.UpdateTrailByID(trailID, name, startLatitude, startLongitude, endLatitude, endLongitude, shelterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update trail"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Trail updated successfully"})
}

func (s *TrailManagerServiceHTTPHandler) GetClosestTrailHandler(c *gin.Context) {
	longitudeStr := c.Query("longitude")
	latitudeStr := c.Query("latitude")

	// Convert query parameters to float64
	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude format"})
		return
	}

	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude format"})
		return
	}

	// Assuming you have a method to find the closest trails by coordinates
	closestTrail, err := s.tvc.tm.FindClosestTrail(longitude, latitude)
	if err != nil {
		// Handle possible errors, such as no trails being found
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve closest trail"})
		return
	}

	// Respond with the ID of the closest trail
	c.JSON(http.StatusOK, gin.H{"closestTrailID": closestTrail})
}



func (s *NewTrailManagerServiceHTTPHandler) CreateShelter(c *gin.Context) {

	sid : = c.Query("id")
	name := c.Query("name")
	Longitude := c.Query("longitude")
	Latitude := c.Query("latitude")

	err = s.tvc.repoS.CreateShelter(sid, name, Latitude, Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trail"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Trail created successfully"})
}

func (s *TrailManagerServiceHTTPHandler) UpdateShelter(c *gin.Context) {
	// Parse the trail ID from the URL parameter or query, depending on your routing setup
	shelterID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trail ID format"})
		return
	}

	name := c.Query("name")
	Longitude, err := strconv.ParseFloat(c.Query("longitude"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid  longitude format"})
		return
	}
	Latitude, err := strconv.ParseFloat(c.Query("latitude"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude format"})
		return
	}
	

	// Call the UpdateTrailByID method from TrailRepository
	err = s.tvc.repoS.UpdateShelterByID(shelterID, name, Latitude, Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update shelter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shelter updated successfully"})
}


func (s *TrailManagerServiceHTTPHandler) GetClosestShelterHandler(c *gin.Context) {
	longitudeStr := c.Query("longitude")
	latitudeStr := c.Query("latitude")

	// Convert query parameters to float64
	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude format"})
		return
	}

	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude format"})
		return
	}

	// Assuming you have a method to find the closest trails by coordinates
	closestShelterId, err := s.tvc.tm.getClosestShelterId(longitude, latitude)
	if err != nil {
		// Handle possible errors, such as no trails being found
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve closest trail"})
		return
	}

	// Respond with the ID of the closest trail
	c.JSON(http.StatusOK, gin.H{"closestShelterID": closestShelterId})
}


