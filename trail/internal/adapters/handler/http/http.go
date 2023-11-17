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
	router.POST("/trail_manager/", handler.CreateTrailManager)

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

func (s *TrailManagerServiceHTTPHandler) CreateTrailManager(ctx *gin.Context) {

	var userDataInstance UserData
	if err := ctx.ShouldBindJSON(&userDataInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	id := userDataInstance.WorkoutID

	_, err := s.tvc.CreateTrailManager(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create trail manager"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "trail manager created successfully"})

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

func (s *TrailManagerServiceHTTPHandler) UpdateTrail(ctx *gin.Context) {

	id, _ := parseUUID(ctx, "trail_id")
	name := ctx.Query("trail_name")
	zId, _ := parseUUID(ctx, "zone_id")
	startLongitude, _ := strconv.ParseFloat(ctx.Query("start_longitude"), 64)
	startLatitude, _ := strconv.ParseFloat(ctx.Query("start_latitude"), 64)
	endLongitude, _ := strconv.ParseFloat(ctx.Query("end_longitude"), 64)
	endLatitude, _ := strconv.ParseFloat(ctx.Query("end_latitude"), 64)

	err := s.tvc.UpdateTrail(id, name, zId, startLatitude, startLongitude, endLatitude, endLongitude)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to update trail"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "updated trail"})

}

func (t *TrailManagerServiceHTTPHandler) DeleteTrail(ctx *gin.Context) {
	id, _ := parseUUID(ctx, "trail_id")
	err := t.tvc.CheckTrail(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to delete trail"})
		return
	}
	err = t.tvc.DeleteTrail(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to delete trail"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "deleted trail"})
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
	name := ctx.Query("shelter_name")
	longitudeStr := ctx.Query("longitude")
	latitudeStr := ctx.Query("latitude")
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)
	latitude, _ := strconv.ParseFloat(latitudeStr, 64)
	availability, _ := strconv.ParseBool("shelter_availability")
	tId, _ := parseUUID(ctx, "trail_id")

	err := t.tvc.UpdateShelter(sId, name, tId, availability, longitude, latitude)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to update shelter"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "shelter updated successfully"})
}

func (t *TrailManagerServiceHTTPHandler) DeleteShelter(ctx *gin.Context) {

	id, _ := parseUUID(ctx, "shelter_id")
	err := t.tvc.CheckShelter(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to create shelter"})
		return
	}
	err = t.tvc.DeleteShelter(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to create shelter"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "shelter deleted successfully"})
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

}

func (h *TrailManagerServiceHTTPHandler) UpdateZone(ctx *gin.Context) {
	id, _ := parseUUID(ctx, "zone_id")
	name := ctx.Query("zone_name")
	err := h.tvc.UpdateZone(id, name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to update shelter"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "updated shelter"})
}

func (h *TrailManagerServiceHTTPHandler) DeleteZone(ctx *gin.Context) {
	id, _ := parseUUID(ctx, "zone_id")
	err := h.tvc.CheckZone(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to delete"})
		return
	}
	err = h.tvc.DeleteZone(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to delete"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "deleted zone"})
}
