package httphandler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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

	router.GET("zone/:zone_id/trail", handler.GetClosestTrail)
	router.GET("zone/:zone_id/trail/:trail_id", handler.GetTrailLocationInfo)
	router.POST("/zone/:zone_id/trail", handler.CreateTrail)
	router.PUT("/zone/:zone_id/trail/:trail_id", handler.UpdateTrail)
	router.DELETE("/zone/:zone_id/trail/:trail_id", handler.DeleteTrail)

	// router.GET("/zone/:zone_id/trail/:trail_id/shelter/check_status", handler.CheckShelterStatus)
	router.GET("/zone/:zone_id/trail/:trail_id/shelter/:shelter_id", handler.GetShelterLocationInfo)
	router.GET("/zone/:zone_id/trail/:trail_id/shelter/", handler.GetClosestShelterInfo)

	router.POST("/zone/:zone_id/trail/:trail_id/shelter", handler.CreateShelter)
	router.PUT("/zone/:zone_id/trail/:trail_id/shelter", handler.UpdateShelter)
	router.DELETE("/zone/:zone_id/trail/:trail_id/shelter/:shelter_id", handler.DeleteShelter)

	router.POST("/zone", handler.CreateZone)
	router.PUT("/zone/:zone_id", handler.UpdateZone)
	router.DELETE("/zone/:zone_id", handler.DeleteZone)

}

func parseUUID(ctx *gin.Context, paramName string) (uuid.UUID, error) {
	uuidStr := ctx.Query(paramName)
	uuidValue, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuidValue, nil
}

func (s *TrailManagerServiceHTTPHandler) GetClosestShelterInfo(ctx *gin.Context) {

	zoneIdStr := ctx.Param("zone_id")
	longitudeStr := ctx.Query("longitude")
	latitudeStr := ctx.Query("latitude")
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)
	latitude, _ := strconv.ParseFloat(latitudeStr, 64)

	_, errZ := uuid.Parse(zoneIdStr)
	if errZ != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errZ.Error()})
		return

	}

	trailIdStr := ctx.Param("trail_id")
	_, errT := uuid.Parse(trailIdStr)
	if errT != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errT.Error()})
		return

	}

	shelterIdStr := ctx.Param("shelter_id")
	id, errS := uuid.Parse(shelterIdStr)
	if errS != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errT.Error()})
		return

	}

	sId, distance, availability, time, err := s.tvc.GetClosestShelter(longitude, latitude, time.Now())
	var shelterDataInstance ShelterAvailable
	shelterDataInstance.WorkoutID = id
	shelterDataInstance.ShelterID = sId
	shelterDataInstance.ShelterAvailable = availability
	shelterDataInstance.ShelterAvailable = (sId == uuid.Nil)
	shelterDataInstance.DistanceToShelter = distance
	shelterDataInstance.ShelterCheckTime = time
	fmt.Println(sId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "no shelter found "})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": shelterDataInstance})

}

func (s *TrailManagerServiceHTTPHandler) CreateTrail(ctx *gin.Context) {

	zoneIdStr := ctx.Param("zone_id")
	zId, errZ := uuid.Parse(zoneIdStr)
	if errZ != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errZ.Error()})
		return

	}

	var trailDataInstance TrailDTO
	if err := ctx.ShouldBindJSON(&trailDataInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	name := trailDataInstance.TrailName
	startLongitude := trailDataInstance.StartLongitude
	startLatitude := trailDataInstance.StartLatitude
	endLongitude := trailDataInstance.EndLongitude
	endLatitude := trailDataInstance.EndLatitude

	tId, err := s.tvc.CreateTrail(name, zId, startLatitude, startLongitude, endLatitude, endLongitude)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trail"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "trail created successfully", "trail_id": tId})

}

func (s *TrailManagerServiceHTTPHandler) UpdateTrail(ctx *gin.Context) {

	zoneIdStr := ctx.Param("zone_id")
	zId, errZ := uuid.Parse(zoneIdStr)
	if errZ != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errZ.Error()})
		return

	}

	trailIdStr := ctx.Param("trail_id")
	id, errT := uuid.Parse(trailIdStr)
	if errT != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errT.Error()})
		return

	}

	name := ctx.Query("trail_name")
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

	zoneIdStr := ctx.Param("zone_id")
	_, errZ := uuid.Parse(zoneIdStr)
	if errZ != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errZ.Error()})
		return

	}

	trailIdStr := ctx.Param("trail_id")
	id, errT := uuid.Parse(trailIdStr)
	if errT != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errT.Error()})
		return

	}

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
	zoneIdStr := ctx.Param("zone_id")
	zId, errZ := uuid.Parse(zoneIdStr)
	if errZ != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errZ.Error()})
		return

	}

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

	zoneIdStr := ctx.Param("zone_id")
	_, errZ := uuid.Parse(zoneIdStr)
	if errZ != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errZ.Error()})
		return

	}

	trailIdStr := ctx.Param("trail_id")
	tId, errT := uuid.Parse(trailIdStr)
	if errT != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errT.Error()})
		return

	}

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

	zoneIdStr := ctx.Param("zone_id")
	_, errZ := uuid.Parse(zoneIdStr)
	if errZ != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errZ.Error()})
		return

	}

	trailIdStr := ctx.Param("trail_id")
	tId, errT := uuid.Parse(trailIdStr)
	if errT != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errT.Error()})
		return

	}

	var shelterDataInstance ShelterDTO
	if err := ctx.ShouldBindJSON(&shelterDataInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	name := shelterDataInstance.ShelterName
	longitude := shelterDataInstance.Longitude
	latitude := shelterDataInstance.Latitude

	sId, err := t.tvc.CreateShelter(name, tId, true, longitude, latitude)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to create shelter"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "shelter created successfully", "shelter_id": sId})
}

func (t *TrailManagerServiceHTTPHandler) UpdateShelter(ctx *gin.Context) {

	zoneIdStr := ctx.Param("zone_id")
	_, errZ := uuid.Parse(zoneIdStr)
	if errZ != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errZ.Error()})
		return

	}

	trailIdStr := ctx.Param("trail_id")
	tId, errT := uuid.Parse(trailIdStr)
	if errT != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errT.Error()})
		return

	}

	var shelterDataInstance ShelterDTO
	if err := ctx.ShouldBindJSON(&shelterDataInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	sId := shelterDataInstance.ShelterID
	name := shelterDataInstance.ShelterName
	longitude := shelterDataInstance.Longitude
	latitude := shelterDataInstance.Latitude
	availability := shelterDataInstance.ShelterAvailability

	err := t.tvc.UpdateShelter(sId, name, tId, availability, longitude, latitude)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to update shelter"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "shelter updated successfully"})
}

func (t *TrailManagerServiceHTTPHandler) DeleteShelter(ctx *gin.Context) {

	zoneIdStr := ctx.Param("zone_id")
	_, errZ := uuid.Parse(zoneIdStr)
	if errZ != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errZ.Error()})
		return

	}

	trailIdStr := ctx.Param("trail_id")
	_, errT := uuid.Parse(trailIdStr)
	if errT != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errT.Error()})
		return

	}

	shelterIdStr := ctx.Param("shelter_id")
	id, errS := uuid.Parse(shelterIdStr)
	if errS != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errT.Error()})
		return

	}

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

func (t *TrailManagerServiceHTTPHandler) GetShelterLocationInfo(ctx *gin.Context) {

	zoneIdStr := ctx.Param("zone_id")
	_, errZ := uuid.Parse(zoneIdStr)
	if errZ != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errZ.Error()})
		return

	}

	trailIdStr := ctx.Param("trail_id")
	_, errT := uuid.Parse(trailIdStr)
	if errT != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errT.Error()})
		return

	}

	shelterIdStr := ctx.Param("shelter_id")
	id, errS := uuid.Parse(shelterIdStr)
	if errS != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": errT.Error()})
		return

	}

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
	zoneIdStr := ctx.Param("zone_id")

	id, _ := uuid.Parse(zoneIdStr)
	var zoneDataInstance ZoneDTO
	if err := ctx.ShouldBindJSON(&zoneDataInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	err := h.tvc.UpdateZone(id, zoneDataInstance.ZoneName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to update shelter"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "updated shelter"})
}

func (h *TrailManagerServiceHTTPHandler) DeleteZone(ctx *gin.Context) {

	zoneIdStr := ctx.Param("zone_id")
	id, err := uuid.Parse(zoneIdStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to delete"})
		return
	}
	err = h.tvc.CheckZone(id)
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
