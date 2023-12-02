package httphandler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/core/services"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type ZoneServiceHTTPHandler struct {
	gin *gin.Engine
	tvc *services.ZoneService
}

func NewZoneServiceHTTPHandler(gin *gin.Engine, tmSvc *services.ZoneService) *ZoneServiceHTTPHandler {
	return &ZoneServiceHTTPHandler{
		gin: gin,
		tvc: tmSvc,
	}
}

func (handler *ZoneServiceHTTPHandler) InitRouter() {

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

func (s *ZoneServiceHTTPHandler) GetClosestShelterInfo(ctx *gin.Context) {

	zoneIdStr := ctx.Param("zone_id")
	longitudeStr := ctx.Query("longitude")
	latitudeStr := ctx.Query("latitude")
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)
	latitude, _ := strconv.ParseFloat(latitudeStr, 64)

	_, err := uuid.Parse(zoneIdStr)
	if err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return

	}

	trailIdStr := ctx.Param("trail_id")
	_, err = uuid.Parse(trailIdStr)
	if err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return

	}

	shelterIdStr := ctx.Param("shelter_id")
	var id uuid.UUID
	id, err = uuid.Parse(shelterIdStr)
	if err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
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

// GetClosestTrail godoc
// @Summary Get the closest trail
// @Description Get the closest trail based on given longitude and latitude
// @Tags zone
// @Accept json
// @Produce json
// @Param zone_id path string true "Zone ID"
// @Param longitude query float64 true "Longitude"
// @Param latitude query float64 true "Latitude"
// @Success 200 {object} Trail
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/zone/{zone_id}/trail [get]

func (s *ZoneServiceHTTPHandler) CreateTrail(ctx *gin.Context) {

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

// UpdateTrail godoc
// @Summary Update a trail
// @Description Update details of an existing trail
// @Tags zone
// @Accept json
// @Produce json
// @Param zone_id path string true "Zone ID"
// @Param trail_id path string true "Trail ID"
// @Param trail body TrailDTO true "Updated Trail Data"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/zone/{zone_id}/trail/{trail_id} [put]
func (s *ZoneServiceHTTPHandler) UpdateTrail(ctx *gin.Context) {
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

// DeleteTrail godoc
// @Summary Delete a specific trail
// @Description Delete a trail from the zone by its ID
// @Tags zone
// @Accept json
// @Produce json
// @Param zone_id path string true "Zone ID"
// @Param trail_id path string true "Trail ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/zone/{zone_id}/trail/{trail_id} [delete]
func (t *ZoneServiceHTTPHandler) DeleteTrail(ctx *gin.Context) {

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

// GetClosestTrail godoc
// @Summary Get the closest trail
// @Description Get the closest trail based on given longitude and latitude
// @Tags zone
// @Accept json
// @Produce json
// @Param zone_id path string true "Zone ID"
// @Param longitude query float64 true "Longitude"
// @Param latitude query float64 true "Latitude"
// @Success 200 {object} Trail
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/zone/{zone_id}/trail [get]
func (t *ZoneServiceHTTPHandler) GetClosestTrail(ctx *gin.Context) {

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

// GetTrailLocationInfo godoc
// @Summary Get location information of a specific trail
// @Description Retrieve detailed location information of a trail by its ID
// @Tags zone
// @Accept json
// @Produce json
// @Param zone_id path string true "Zone ID"
// @Param trail_id path string true "Trail ID"
// @Success 200 {object} TrailLocationInfo
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/zone/{zone_id}/trail/{trail_id} [get]
func (t *ZoneServiceHTTPHandler) GetTrailLocationInfo(ctx *gin.Context) {

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
// CreateShelter godoc
// @Summary Create a shelter
// @Description Create a new shelter associated with a trail
// @Tags zone
// @Accept json
// @Produce json
// @Param zone_id path string true "Zone ID"
// @Param trail_id path string true "Trail ID"
// @Param shelter body ShelterDTO true "Shelter Data"
// @Success 201 {object} ShelterCreationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/zone/{zone_id}/trail/{trail_id}/shelter [post]

func (t *ZoneServiceHTTPHandler) CreateShelter(ctx *gin.Context) {

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

// UpdateShelter godoc
// @Summary Update a shelter
// @Description Update details of an existing shelter
// @Tags zone
// @Accept json
// @Produce json
// @Param zone_id path string true "Zone ID"
// @Param trail_id path string true "Trail ID"
// @Param shelter body ShelterDTO true "Updated Shelter Data"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/zone/{zone_id}/trail/{trail_id}/shelter [put]

func (t *ZoneServiceHTTPHandler) UpdateShelter(ctx *gin.Context) {

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

// DeleteShelter godoc
// @Summary Delete a shelter
// @Description Delete an existing shelter
// @Tags zone
// @Accept json
// @Produce json
// @Param zone_id path string true "Zone ID"
// @Param trail_id path string true "Trail ID"
// @Param shelter_id path string true "Shelter ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/zone/{zone_id}/trail/{trail_id}/shelter/{shelter_id} [delete]

func (t *ZoneServiceHTTPHandler) DeleteShelter(ctx *gin.Context) {

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

// GetShelterLocationInfo godoc
// @Summary Get location information of a specific shelter
// @Description Retrieve detailed location information of a shelter by its ID
// @Tags zone
// @Accept json
// @Produce json
// @Param zone_id path string true "Zone ID"
// @Param trail_id path string true "Trail ID"
// @Param shelter_id path string true "Shelter ID"
// @Success 200 {object} ShelterLocationInfo
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/zone/{zone_id}/trail/{trail_id}/shelter/{shelter_id} [get]
func (t *ZoneServiceHTTPHandler) GetShelterLocationInfo(ctx *gin.Context) {

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

// CreateZone godoc
// @Summary Create a zone
// @Description Create a new zone
// @Tags zone
// @Accept json
// @Produce json
// @Param zone body ZoneDTO true "Zone Data"
// @Success 201 {object} ZoneCreationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/zone [post]

func (h *ZoneServiceHTTPHandler) CreateZone(ctx *gin.Context) {
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

// UpdateZone godoc
// @Summary Update a zone
// @Description Update details of an existing zone
// @Tags zone
// @Accept json
// @Produce json
// @Param zone_id path string true "Zone ID"
// @Param zone body ZoneDTO true "Updated Zone Data"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/zone/{zone_id} [put]

func (h *ZoneServiceHTTPHandler) UpdateZone(ctx *gin.Context) {
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

// DeleteZone godoc
// @Summary Delete a zone
// @Description Delete an existing zone
// @Tags zone
// @Accept json
// @Produce json
// @Param zone_id path string true "Zone ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/zone/{zone_id} [delete]

func (h *ZoneServiceHTTPHandler) DeleteZone(ctx *gin.Context) {

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
