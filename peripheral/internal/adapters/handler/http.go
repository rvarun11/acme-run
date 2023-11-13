package handler

import (
	"net/http"

	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/services"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	svc services.PeripheralService
}

func PeripheralServiceHTTPHandler(PeripheralService services.PeripheralService) *HTTPHandler {
	return &HTTPHandler{
		svc: PeripheralService,
	}
}

func (handler *PeripheralServiceHTTPHandler) InitRouter() {

	router := handler.gin.Group("/api/v1")

	router.POST("/peripheral", handler.BindPeripheralToData)
	router.PUT("/peripheral/hrm_status", handler.SetHRMStatus)
	router.GET("/peripheral/hrm_staus", handler.GetHRMStatus)
	router.GET("/peripheral/hrm_reading", handler.GetHRMReading)
	router.PUT("/peripheral/hrm_reading", handler.SetHRMReading)
	router.PUT("/peripheral/geo_status", handler.SetGeoStatus)
	router.GET("/peripheral/geo_status", handler.GetGeoStatus)
	router.GET("/peripheral/geo_reading", handler.GetGeoReading)
	router.PUT("/peripheral/geo_reading", handler.SetGeoReading)

}

func (h *HTTPHandler) BindPeripheralToData(ctx *gin.Context) {

	var bindDataInstance BindPeripheralData
	if err := ctx.ShouldBindJSON(&bindDataInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}

	if bindDataInstance.Connect {
		go handler.rabbitMQHandler.SendLastHR(data.WorkoutID)
		h.svc.ConnectPeripheral(bindDataInstance.WorkoutId, bindDataInstance.HRMId)
	} else {
		h.svc.DisconnectPeripheral(bindDataInstance.WorkoutId)
	}
}

func (h *HTTPHandler) GetHRMReading(ctx *gin.Context) {
	wId, err := parseUUID(ctx, "workout_id")
	ctx.JSON(http.StatusOK, h.svc.GetHRMReading(wId))
}

func (h *HTTPHandler) GetHRMStatus(ctx *gin.Context) {

	wId := ctx.Query("workout_id")
	ctx.JSON(http.StatusOK, h.svc.GetHRMDevStatus(wId))
}

func (h *HTTPHandler) SetHRMStatus(ctx *gin.Context) {

	wId := ctx.Query("workout_id")
	code := ctx.Query("code")
	ctx.JSON(http.StatusOK, h.svc.SetHRMDevStatus(code))
}

func (h *HTTPHandler) SetHRMReading(ctx *gin.Context) {
	wId, err := parseUUID(ctx, "workout_id")
	rate := ctx.Query("current_reading")
	ctx.JSON(http.StatusOK, h.svc.SetAverageHRate(wId, rate))
}

func (h *HTTPHandler) GetGeoStatus(ctx *gin.Context) {

	wId := ctx.Query("workout_id")
	ctx.JSON(http.StatusOK, h.svc.GetGeoDevStatus(wId))
}

func (h *HTTPHandler) SetGeoStatus(ctx *gin.Context) {

	wId := ctx.Query("workout_id")
	code := ctx.Query("code")
	ctx.JSON(http.StatusOK, h.svc.GetGeoDevStatus(code))
}

func (h *HTTPHandler) SetGeoReading(ctx *gin.Context) {

	latitude := ctx.Query("latitude")
	longitude := ctx.Query("longitude")
	wId := ctx.Query("workout_id")

	var tempLastLoc LastLocation
	tempLastLoc.WorkoutID = wId
	tempLastLoc.TimeOfLocation = time.Time()
	h.svc.SetGeoLocation(wId, longitude, latitude)

	if err := c.ShouldBindJSON(&tempLastLoc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Now trigger the RabbitMQHandler to send the updated location
	go handler.rabbitMQHandler.SendLastLocation(data.WorkoutID)
	c.JSON(http.StatusOK, gin.H{"message": "Geo reading set and location sent"})

	// ctx.JSON(http.StatusOK, h.svc.SetGeoLocation(wId, longitude, latitude))
}

func (h *HTTPHandler) GetGeoReading(ctx *gin.Context) {
	wId := ctx.Query("workout_id")
	ctx.JSON(http.StatusOK, h.svc.GetGeoLocation(wId))
}
