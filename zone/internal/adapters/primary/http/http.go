package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/core/services"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type ZoneHandler struct {
	gin *gin.Engine
	tvc *services.ZoneService
}

func NewZoneHandler(gin *gin.Engine, tmSvc *services.ZoneService) *ZoneHandler {
	return &ZoneHandler{
		gin: gin,
		tvc: tmSvc,
	}
}

func (handler *ZoneHandler) InitRouter() {

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
	router.PUT("/zone/:zone_id/trail/:trail_id/shelter/:shelter_id", handler.UpdateShelter)
	router.DELETE("/zone/:zone_id/trail/:trail_id/shelter/:shelter_id", handler.DeleteShelter)

	router.POST("/zone", handler.CreateZone)
	router.PUT("/zone/:zone_id", handler.UpdateZone)
	router.DELETE("/zone/:zone_id", handler.DeleteZone)

}

// GetClosestShelterInfo
//
//	@Summary		Get the closest shelter information
//	@Description	Retrieve the closest shelter information based on current longitude and latitude
//	@Tags			zone
//	@Accept			json
//	@Router			/api/v1/zone/{zone_id}/trail/{trail_id}/shelter [get]
func (s *ZoneHandler) GetClosestShelterInfo(ctx *gin.Context) {

	zoneIdStr := ctx.Param("zone_id")
	longitudeStr := ctx.Query("longitude")
	latitudeStr := ctx.Query("latitude")
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)
	latitude, _ := strconv.ParseFloat(latitudeStr, 64)

	_, err := uuid.Parse(zoneIdStr)
	if err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid zone id"})
		return

	}

	trailIdStr := ctx.Param("trail_id")
	_, err = uuid.Parse(trailIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid zone id"})
		return

	}

	shelterIdStr := ctx.Param("shelter_id")
	var id uuid.UUID
	id, err = uuid.Parse(shelterIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid trail id"})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "shelter not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"distance_to_shelter": shelterDataInstance})

}

// CreateTrail
//
//	@Summary		Create a trail
//	@Description	Create a trail based on given longitude and latitude within a zone
//	@Tags			zone
//	@Accept			json
//	@Produce		json
//	@Param			zone_id	path		string				true	"Zone ID"
//	@Param			trail	body		TrailDTO			true	"Trail Data"
//	@Success		201		{object}	map[string]string	"status: success, message: trail created successfully, trail_id: UUID"
//	@Failure		400		{object}	map[string]string	"status: error, message: failed to create trail"
//	@Failure		500		{object}	map[string]string	"status: error, message: Internal Server Error"
//	@Router			/api/v1/zone/{zone_id}/trail [post]
func (s *ZoneHandler) CreateTrail(ctx *gin.Context) {

	zoneIdStr := ctx.Param("zone_id")
	zId, errZ := uuid.Parse(zoneIdStr)
	if errZ != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid zone id"})
		return

	}

	var trailDataInstance TrailDTO
	if err := ctx.ShouldBindJSON(&trailDataInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request paramenters"})
		return
	}

	name := trailDataInstance.TrailName
	startLongitude := trailDataInstance.StartLongitude
	startLatitude := trailDataInstance.StartLatitude
	endLongitude := trailDataInstance.EndLongitude
	endLatitude := trailDataInstance.EndLatitude

	tId, err := s.tvc.CreateTrail(name, zId, startLatitude, startLongitude, endLatitude, endLongitude)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create trail, something went wrong"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "trail created successfully", "trail_id": tId})

}

// UpdateTrail
//
//	@Summary		Update a trail
//	@Description	Update details of an existing trail within a zone
//	@Tags			zone
//	@Accept			json
//	@Produce		json
//	@Param			zone_id		path		string				true	"Zone ID"
//	@Param			trail_id	path		string				true	"Trail ID"
//	@Param			trail		body		TrailDTO			true	"Updated Trail Data"
//	@Success		200			{object}	map[string]string	"status: success, message: updated trail"
//	@Failure		400			{object}	map[string]string	"status: error, message: failed to update trail"
//	@Failure		500			{object}	map[string]string	"status: error, message: Internal Server Error"
//	@Router			/api/v1/zone/{zone_id}/trail/{trail_id} [put]
func (s *ZoneHandler) UpdateTrail(ctx *gin.Context) {
	zoneIdStr := ctx.Param("zone_id")
	zId, errZ := uuid.Parse(zoneIdStr)
	if errZ != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid zone id"})
		return

	}

	trailIdStr := ctx.Param("trail_id")
	id, errT := uuid.Parse(trailIdStr)
	if errT != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid trail id"})
		return

	}

	name := ctx.Query("trail_name")
	startLongitude, _ := strconv.ParseFloat(ctx.Query("start_longitude"), 64)
	startLatitude, _ := strconv.ParseFloat(ctx.Query("start_latitude"), 64)
	endLongitude, _ := strconv.ParseFloat(ctx.Query("end_longitude"), 64)
	endLatitude, _ := strconv.ParseFloat(ctx.Query("end_latitude"), 64)

	err := s.tvc.UpdateTrail(id, name, zId, startLatitude, startLongitude, endLatitude, endLongitude)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update trail, something went wrong"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "updated trail successfully"})

}

// DeleteTrail
//
//	@Summary		Delete a specific trail
//	@Description	Delete a trail from the zone by its ID
//	@Tags			zone
//	@Accept			json
//	@Produce		json
//	@Param			zone_id		path		string				true	"Zone ID"
//	@Param			trail_id	path		string				true	"Trail ID"
//	@Success		200			{object}	map[string]string	"status: success, message: deleted trail"
//	@Failure		400			{object}	map[string]string	"status: error, message: failed to delete trail"
//	@Failure		500			{object}	map[string]string	"status: error, message: Internal Server Error"
//	@Router			/api/v1/zone/{zone_id}/trail/{trail_id} [delete]
func (t *ZoneHandler) DeleteTrail(ctx *gin.Context) {

	zoneIdStr := ctx.Param("zone_id")
	_, errZ := uuid.Parse(zoneIdStr)
	if errZ != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid zone id"})
		return

	}

	trailIdStr := ctx.Param("trail_id")
	id, errT := uuid.Parse(trailIdStr)
	if errT != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid trail id"})
		return

	}

	err := t.tvc.CheckTrail(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no trail found"})
		return
	}
	err = t.tvc.DeleteTrail(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete trail, something went wrong"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "deleted trail successfully"})
}

// GetClosestTrail
//
//	@Summary		Get the closest trail
//	@Description	Get the closest trail based on given longitude and latitude in a specific zone
//	@Tags			zone
//	@Accept			json
//	@Produce		json
//	@Param			zone_id		path		string				true	"Zone ID"
//	@Param			longitude	query		float64				true	"Longitude"
//	@Param			latitude	query		float64				true	"Latitude"
//	@Success		200			{object}	map[string]string	"status: success, message: closest trail, trail_id: UUID"
//	@Failure		400			{object}	map[string]string	"status: error, message: failed to retrieve closest trail"
//	@Failure		500			{object}	map[string]string	"status: error, message: Internal Server Error"
//	@Router			/api/v1/zone/{zone_id}/trail [get]
func (t *ZoneHandler) GetClosestTrail(ctx *gin.Context) {

	longitudeStr := ctx.Query("longitude")
	latitudeStr := ctx.Query("latitude")
	zoneIdStr := ctx.Param("zone_id")
	zId, errZ := uuid.Parse(zoneIdStr)
	if errZ != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid zone id"})
		return

	}

	// Convert query parameters to float64
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)

	latitude, _ := strconv.ParseFloat(latitudeStr, 64)

	// Assuming you have a method to find the closest trails by coordinates
	closestTrail, err := t.tvc.GetClosestTrail(zId, longitude, latitude)
	if err != nil || closestTrail == uuid.Nil {
		// Handle possible errors, such as no trails being found
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve closest trail, something went wrong"})
		return
	}

	// Respond with the ID of the closest trail

	ctx.JSON(http.StatusOK, gin.H{"trail_id": closestTrail})
}

// GetTrailLocationInfo
//
//	@Summary		Get location information of a specific trail
//	@Description	Retrieve detailed location information of a trail by its ID
//	@Tags			zone
//	@Accept			json
//	@Produce		json
//	@Param			zone_id		path		string				true	"Zone ID"
//	@Param			trail_id	path		string				true	"Trail ID"
//	@Success		200			{string}	string				"UUID of the closest trail"
//	@Failure		400			{object}	map[string]string	"status: error, message: Failed"
//	@Failure		500			{object}	map[string]string	"status: error, message: Failed"
//	@Router			/api/v1/zone/{zone_id}/trail/{trail_id} [get]
func (t *ZoneHandler) GetTrailLocationInfo(ctx *gin.Context) {

	zoneIdStr := ctx.Param("zone_id")
	_, errZ := uuid.Parse(zoneIdStr)
	if errZ != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid zone id"})
		return

	}

	trailIdStr := ctx.Param("trail_id")
	tId, errT := uuid.Parse(trailIdStr)
	if errT != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid trail id"})

	}

	// Assuming you have a method to find the closest trails by coordinates
	trail, err := t.tvc.GetTrailByID(tId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve trail info, something went wrong"})
		return
	}

	// Respond with the ID of the closest trail
	ctx.JSON(http.StatusOK, gin.H{"start_longitude": trail.StartLongitude, "start_latitude": trail.StartLatitude,
		"end_longitude": trail.EndLongitude, "end_latitude": trail.EndLatitude})
}

// CreateShelter
//
//	@Summary		Create a shelter
//	@Description	Create a new shelter associated with a trail in a zone
//	@Tags			zone
//	@Accept			json
//	@Produce		json
//	@Param			zone_id		path		string				true	"Zone ID"
//	@Param			trail_id	path		string				true	"Trail ID"
//	@Param			shelter		body		ShelterDTO			true	"Shelter Data"
//	@Success		201			{object}	map[string]string	"status: success, message: shelter created successfully, shelter_id: UUID"
//	@Failure		400			{object}	map[string]string	"status: error, message: failed to create shelter"
//	@Failure		500			{object}	map[string]string	"status: error, message: Internal Server Error"
//	@Router			/api/v1/zone/{zone_id}/trail/{trail_id}/shelter [post]
func (t *ZoneHandler) CreateShelter(ctx *gin.Context) {

	zoneIdStr := ctx.Param("zone_id")
	_, errZ := uuid.Parse(zoneIdStr)
	if errZ != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid zone id"})
		return

	}

	trailIdStr := ctx.Param("trail_id")
	tId, errT := uuid.Parse(trailIdStr)
	if errT != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid trail id"})
		return

	}

	var shelterDataInstance ShelterDTO
	if err := ctx.ShouldBindJSON(&shelterDataInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request parameters"})
		return
	}
	name := shelterDataInstance.ShelterName
	longitude := shelterDataInstance.Longitude
	latitude := shelterDataInstance.Latitude

	sId, err := t.tvc.CreateShelter(name, tId, true, longitude, latitude)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create shelter, something went wrong"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "shelter created successfully", "shelter_id": sId})
}

// UpdateShelter
//
//	@Summary		Update a shelter
//	@Description	Update details of an existing shelter in a trail
//	@Tags			zone
//	@Accept			json
//	@Produce		json
//	@Param			zone_id		path		string				true	"Zone ID"
//	@Param			trail_id	path		string				true	"Trail ID"
//	@Param			shelter		body		ShelterDTO			true	"Updated Shelter Data"
//	@Success		200			{object}	map[string]string	"status: success, message: shelter updated successfully"
//	@Failure		400			{object}	map[string]string	"status: error, message: failed to update shelter"
//	@Failure		500			{object}	map[string]string	"status: error, message: Internal Server Error"
//	@Router			/api/v1/zone/{zone_id}/trail/{trail_id}/shelter/{shelter_id} [put]
func (t *ZoneHandler) UpdateShelter(ctx *gin.Context) {

	zoneIdStr := ctx.Param("zone_id")
	_, errZ := uuid.Parse(zoneIdStr)
	if errZ != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid zone id"})
		return

	}

	trailIdStr := ctx.Param("trail_id")
	tId, errT := uuid.Parse(trailIdStr)
	if errT != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid trail id"})
		return

	}

	shelterIdStr := ctx.Param("shelter_id")
	sId, errS := uuid.Parse(shelterIdStr)
	if errS != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid shelter id"})
		return

	}

	var shelterDataInstance ShelterDTO
	if err := ctx.ShouldBindJSON(&shelterDataInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request parameters"})
		return
	}

	name := shelterDataInstance.ShelterName
	longitude := shelterDataInstance.Longitude
	latitude := shelterDataInstance.Latitude
	availability := shelterDataInstance.ShelterAvailability

	err := t.tvc.UpdateShelter(sId, name, tId, availability, longitude, latitude)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update shelter, something went wrong"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "shelter updated successfully"})
}

// DeleteShelter
//
//	@Summary		Delete a shelter
//	@Description	Delete an existing shelter from a trail in a zone
//	@Tags			zone
//	@Accept			json
//	@Produce		json
//	@Param			zone_id		path		string				true	"Zone ID"
//	@Param			trail_id	path		string				true	"Trail ID"
//	@Param			shelter_id	path		string				true	"Shelter ID"
//	@Success		200			{object}	map[string]string	"status: success, message: shelter deleted successfully"
//	@Failure		400			{object}	map[string]string	"status: error, message: failed to delete shelter"
//	@Failure		500			{object}	map[string]string	"status: error, message: Internal Server Error"
//	@Router			/api/v1/zone/{zone_id}/trail/{trail_id}/shelter/{shelter_id} [delete]
func (t *ZoneHandler) DeleteShelter(ctx *gin.Context) {

	zoneIdStr := ctx.Param("zone_id")
	_, errZ := uuid.Parse(zoneIdStr)
	if errZ != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid zone id"})
		return

	}

	trailIdStr := ctx.Param("trail_id")
	_, errT := uuid.Parse(trailIdStr)
	if errT != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid trail id"})
		return

	}

	shelterIdStr := ctx.Param("shelter_id")
	id, errS := uuid.Parse(shelterIdStr)
	if errS != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid shelter id"})
		return

	}

	err := t.tvc.CheckShelter(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find shelter, something went wrong"})
		return
	}
	err = t.tvc.DeleteShelter(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete shelter, something went wrong"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "shelter deleted successfully"})
}

// GetShelterLocationInfo
//
//	@Summary		Get location information of a specific shelter
//	@Description	Retrieve detailed location information of a shelter by its ID in a trail
//	@Tags			zone
//	@Accept			json
//	@Produce		json
//	@Param			zone_id		path		string				true	"Zone ID"
//	@Param			trail_id	path		string				true	"Trail ID"
//	@Param			shelter_id	path		string				true	"Shelter ID"
//	@Success		200			{object}	map[string]float64	"status: success, message: shelter location information"
//	@Failure		400			{object}	map[string]string	"status: error, message: failed to retrieve shelter info"
//	@Failure		500			{object}	map[string]string	"status: error, message: Internal Server Error"
//	@Router			/api/v1/zone/{zone_id}/trail/{trail_id}/shelter/{shelter_id} [get]
func (t *ZoneHandler) GetShelterLocationInfo(ctx *gin.Context) {

	zoneIdStr := ctx.Param("zone_id")
	_, errZ := uuid.Parse(zoneIdStr)
	if errZ != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid zone id"})
		return

	}

	trailIdStr := ctx.Param("trail_id")
	_, errT := uuid.Parse(trailIdStr)
	if errT != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid trail id"})
		return

	}

	shelterIdStr := ctx.Param("shelter_id")
	id, errS := uuid.Parse(shelterIdStr)
	if errS != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid shelter id"})
		return

	}

	// Assuming you have a method to find the closest trails by coordinates
	shelter, err := t.tvc.GetShelterByID(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to retrieve shelter info, something went wrong"})
		return
	}

	// Respond with the ID of the closest trail
	ctx.JSON(http.StatusOK, gin.H{"longitude": shelter.Longitude, "latitude": shelter.Latitude})
}

// CreateZone
//
//	@Summary		Create a zone
//	@Description	Create a new zone
//	@Tags			zone
//	@Accept			json
//	@Produce		json
//	@Param			zone	body		ZoneDTO				true	"Zone Data"
//	@Success		201		string		"message":			"zone created"
//	@Failure		400		{object}	map[string]string	"status: error, message: Failed"
//	@Failure		500		{object}	map[string]string	"status: error, message: Failed"
//	@Router			/api/v1/zone [post]
func (h *ZoneHandler) CreateZone(ctx *gin.Context) {
	var zoneDataInstance ZoneDTO
	if err := ctx.ShouldBindJSON(&zoneDataInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request parameters"})
		return
	}
	zId, err := h.tvc.CreateZone(zoneDataInstance.ZoneName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error while creating the zone, something went wrong"})
	} else {
		zIdString := zId.String()

		ctx.JSON(http.StatusCreated, gin.H{"message": "zone created", "zone_id": zIdString})
	}

}

// UpdateZone
//
//	@Summary		Update a zone
//	@Description	Update details of an existing zone
//	@Tags			zone
//	@Accept			json
//	@Produce		json
//	@Param			zone_id	path		string				true	"Zone ID"
//	@Param			zone	body		ZoneDTO				true	"Updated Zone Data"
//	@Success		200		{object}	map[string]string	"status: error, message: success"
//	@Failure		400		{object}	map[string]string	"status: error, message: Failed"
//	@Failure		500		{object}	map[string]string	"status: error, message: Failed"
//	@Router			/api/v1/zone/{zone_id} [put]
func (h *ZoneHandler) UpdateZone(ctx *gin.Context) {
	zoneIdStr := ctx.Param("zone_id")

	id, _ := uuid.Parse(zoneIdStr)
	var zoneDataInstance ZoneDTO
	if err := ctx.ShouldBindJSON(&zoneDataInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request parameters"})
		return
	}
	err := h.tvc.UpdateZone(id, zoneDataInstance.ZoneName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update shelter, something went wrong"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "updated shelter"})
}

// DeleteZone
//
//	@Summary		Delete a zone
//	@Description	Delete an existing zone
//	@Tags			zone
//	@Accept			json
//	@Produce		json
//	@Param			zone_id	path		string				true	"Zone ID"
//	@Success		200		{object}	map[string]string	"status: error, message: success"
//	@Failure		400		{object}	map[string]string	"status: error, message: Failed"
//	@Failure		500		{object}	map[string]string	"status: error, message: Failed"
//	@Router			/api/v1/zone/{zone_id} [delete]
func (h *ZoneHandler) DeleteZone(ctx *gin.Context) {

	zoneIdStr := ctx.Param("zone_id")
	id, err := uuid.Parse(zoneIdStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invlaid zone id"})
		return
	}
	err = h.tvc.CheckZone(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find zone, something went wrong"})
		return
	}
	err = h.tvc.DeleteZone(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete, something went wrong"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "deleted zone"})
}
