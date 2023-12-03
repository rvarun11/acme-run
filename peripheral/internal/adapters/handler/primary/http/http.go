package httphandler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	rabbitmqhandler "github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/adapters/handler/primary/rabbitmq"
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/services"
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/log"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func parseUUID(ctx *gin.Context, paramName string) (uuid.UUID, error) {
	uuidStr := ctx.Query(paramName)
	uuidValue, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuidValue, nil
}

func randomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

type HTTPHandler struct {
	gin             *gin.Engine
	svc             *services.PeripheralService
	rabbitMQHandler *rabbitmqhandler.RabbitMQHandler
	cancelF         context.CancelFunc
	bCtx            context.Context
	hLiveCount      int
}

func NewPeripheralServiceHTTPHandler(gin *gin.Engine, PeripheralService *services.PeripheralService, rabbitMQHandler *rabbitmqhandler.RabbitMQHandler) *HTTPHandler {
	return &HTTPHandler{
		gin:             gin,
		svc:             PeripheralService,
		rabbitMQHandler: rabbitMQHandler,
	}
}

func (handler *HTTPHandler) InitRouter() {

	router := handler.gin.Group("/api/v1")

	router.POST("/peripheral", handler.BindPeripheralToData)
	router.PUT("/peripheral", handler.UnbindPeripheralToData)
	router.PUT("/hrm_status", handler.SetHRMStatus)

	// HRM
	router.POST("/peripheral/hrm", handler.connectHRM)
	router.PUT("/peripheral/hrm", handler.disconnectHRM)
	router.GET("/peripheral/hrm", handler.getHRMReading)

	router.GET("/hrm_status", handler.GetHRMStatus)
	router.PUT("/hrm_reading", handler.SetHRMReading)
	router.PUT("/geo_status", handler.SetGeoStatus)
	router.GET("/geo_status", handler.GetGeoStatus)
	router.GET("/geo_reading", handler.GetGeoReading)
	router.PUT("/geo_reading", handler.SetGeoReading)

}

// ConnectHRM connects to a Heart Rate Monitor (HRM) device.
//
//	@Summary	Connect to HRM device
//	@Tags		peripheral
//	@ID			connect-hrm
//	@Accept		json
//	@Produce	json
//	@Param		connectData	body		BindPeripheralData	true	"Connect Peripheral Data"
//	@Success	200			{object}	map[string]bool		"connect to hrm success: true"
//	@Failure	400			{object}	map[string]string	"status: error, message: Invalid request"
//	@Failure	500			{object}	map[string]string	"error: Cannot connect to HRM"
//	@Router		/api/v1/peripheral/hrm [post]
func (h *HTTPHandler) connectHRM(ctx *gin.Context) {
	var req BindPeripheralData
	var err error
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error", "message": "invalid request",
		})
		return
	}

	if !h.svc.CheckStatusByHRMId(req.HRMId) {
		err = h.svc.CreatePeripheral(req.PlayerID, req.HRMId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	err2 := h.svc.SetHRMDevStatusByHRMId(req.HRMId, true)
	if err2 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot connect to hrm",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"connect to hrm success": true})
}

// DisconnectHRM disconnects a Heart Rate Monitor (HRM) device.
//
//	@Summary	Disconnect HRM device
//	@Tags		peripheral
//	@ID			disconnect-hrm
//	@Accept		json
//	@Produce	json
//	@Param		disconnectData	body		BindPeripheralData	true	"Disconnect Peripheral Data"
//	@Success	200				{object}	map[string]string	"status: success, message: Disconnected HRM"
//	@Failure	400				{object}	map[string]string	"status: error, message: Invalid request"
//	@Failure	500				{object}	map[string]string	"status: error, message: Cannot disconnect HRM"
//	@Router		/api/v1/peripheral/hrm [put]
func (h *HTTPHandler) disconnectHRM(ctx *gin.Context) {

	var req BindPeripheralData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error", "message": "invalid request",
		})
		return
	}
	if !req.HRMConnect {
		err2 := h.svc.SetHRMDevStatusByHRMId(req.HRMId, false)
		if err2 != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "error", "message": "cannot disconnect hrm",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "disconnected hrm"})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error", "message": "cannot disconnect hrm, wrong value of connect",
		})
		return
	}
}

// BindPeripheralToData connects to an HRM device and binds it to a workout
//
//	@Summary	Connect to HRM device and bind it to a workout
//	@Tags		peripheral
//	@ID			bind-peripheral
//	@Accept		json
//	@Produce	json
//	@Param		bindData	body		BindPeripheralData	true	"Bind Peripheral Data"
//	@Success	200			{object}	map[string]string	"status: success, message: Binding workout done"
//	@Failure	400			{object}	map[string]string	"status: error, message: error message"
//	@Router		/api/v1/peripheral [post]
func (h *HTTPHandler) BindPeripheralToData(ctx *gin.Context) {

	var bindDataInstance BindPeripheralData
	if err := ctx.ShouldBindJSON(&bindDataInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	err := h.svc.BindPeripheral(bindDataInstance.PlayerID, bindDataInstance.WorkoutID, bindDataInstance.HRMId, bindDataInstance.HRMConnect, bindDataInstance.SendLiveLocationToTrailManager)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to Bind workout",
		})
		return
	} else {

		h.hLiveCount += 1
		h.bCtx, h.cancelF = context.WithCancel(context.Background())

		longitudeStart, latitudeStart, longitudeEnd, latitudeEnd, err := h.svc.GetTrailLocationInfo(bindDataInstance.TrailOfWorkout)
		if err != nil {
			log.Error("Peripheral: failed to get trail location, using default info now", zap.Error(err))
			longitudeStart = -79.919390
			latitudeStart = 43.257715
			longitudeEnd = 43.258012
			latitudeEnd = -79.910866
		}

		h.svc.SetLiveStatus(bindDataInstance.WorkoutID, true)
		h.StartBackgroundMockReading(ctx, h.bCtx, bindDataInstance.WorkoutID, bindDataInstance.HRMId, longitudeStart, latitudeStart, longitudeEnd, latitudeEnd)
		log.Info("Peripheral: unbind status success")
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Binding workout done",
		})
	}
}

// UnbindPeripheralToData unbinds peripheral data from a workout.
//
//	@Summary	Unbind peripheral data from a workout
//	@Tags		peripheral
//	@ID			unbind-peripheral
//	@Accept		json
//	@Produce	json
//	@Param		unbindData	body		UnbindPeripheralData	true	"Unbind Peripheral Data"
//	@Success	200			{object}	map[string]string		"status: success, message: Unbind the data"
//	@Failure	400			{object}	map[string]string		"status: error, message: Invalid request"
//	@Failure	500			{object}	map[string]string		"status: error, message: Failed to unbind"
//	@Router		/api/v1/peripheral [put]
func (h *HTTPHandler) UnbindPeripheralToData(ctx *gin.Context) {
	var req UnbindPeripheralData
	var err error
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request",
		})
		return
	}

	err = h.svc.SetLiveStatus(req.WorkoutID, false)
	if err != nil {
		log.Error("Peripheral: failed to set live status of publising ", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "failed to unbind",
		})
		return
	}

	h.hLiveCount -= 1
	if h.hLiveCount == 0 {
		h.cancelF()
	}
	log.Info("Peripheral: unbind status success")
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Unbind the data"})
}

// getHRMReading retrieves Heart Rate Monitor (HRM) reading data.
//
//	@Summary	Get HRM reading data
//	@Tags		peripheral
//	@ID			get-hrm-reading
//	@Accept		json
//	@Produce	json
//	@Param		workout_id	path		string				true	"Workout ID"	format(uuid)
//	@Param		type		query		string				true	"Type of HRM reading (avg/normal)"
//	@Success	200			{object}	LastHR				"HRM reading data"
//	@Failure	400			{object}	map[string]string	"status: error, message: Invalid request"
//	@Failure	500			{object}	map[string]string	"status: error, message: Reading from device failure"
//	@Router		/api/v1/peripheral/hrm/{workout_id} [get]
func (h *HTTPHandler) getHRMReading(ctx *gin.Context) {
	wId, err1 := parseUUID(ctx, "workout_id")
	if err1 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "get hrm reading invalid request"})
		return
	}

	hrType := ctx.Query("type")
	if hrType == "avg" {
		// TODO: This should be returning as per workout
		var tLoc LastHR
		var err error
		tLoc.HRMID, tLoc.TimeOfLocation, tLoc.HeartRate, err = h.svc.GetHRMAvgReading(wId)
		if err != nil {
			log.Error("Peripheral: failed to read from device failure ", zap.Error(err))
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "error", "message": "reading from device failure",
			})
			return
		}
		avgRate := AverageHeartRate{}
		avgRate.WorkoutID = wId
		avgRate.AverageHeartRate = uint8(tLoc.HeartRate)

		jsonData, err := json.Marshal(avgRate)
		if err != nil {
			log.Error("Peripheral: failed to read from device,failed to marshal json ", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "error", "message": "could not marshal JSON for average heart rate",
			})
			return
		}
		// ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": jsonData})
		ctx.Writer.Header().Set("Content-Type", "application/json")
		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Writer.Write(jsonData)
	} else if hrType == "normal" {
		var tLoc LastHR
		var err error
		tLoc.HRMID, tLoc.TimeOfLocation, tLoc.HeartRate, err = h.svc.GetHRMReading(wId)
		if err != nil {
			log.Error("Peripheral: failed to read from device failure ", zap.Error(err))
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "error", "message": "reading from device failure",
			})
			return
		}
		log.Info("Peripheral: get hrm reading success")
		ctx.JSON(http.StatusOK, gin.H{"reading": tLoc})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error", "message": "reading from device failure",
		})
	}
}

// GetHRMStatus retrieves the Heart Rate Monitor (HRM) device status for a workout.
//
//	@Summary	Get HRM device status
//	@Tags		peripheral
//	@ID			get-hrm-status
//	@Accept		json
//	@Produce	json
//	@Param		workout_id	path		string				true	"Workout ID"	format(uuid)
//	@Success	200			{object}	map[string]bool		"status: success, message: HRM device status"
//	@Failure	400			{object}	map[string]string	"status: error, message: Failed to get HRM status"
//	@Router		/api/v1/hrm_status [get]
func (h *HTTPHandler) GetHRMStatus(ctx *gin.Context) {

	wId, err := parseUUID(ctx, "workout_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "failed to get hrm status"})
		return
	}

	tStatus, err := h.svc.GetHRMDevStatus(wId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Fail to get hrm status"})
		return
	}
	log.Info("Peripheral: get geo device status success")
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": tStatus})

}

// SetHRMStatus sets the Heart Rate Monitor (HRM) device status for a workout.
//
//	@Summary	Set HRM device status
//	@Tags		peripheral
//	@ID			set-hrm-status
//	@Accept		json
//	@Produce	json
//	@Param		workout_id	path		string				true	"Workout ID"	format(uuid)
//	@Param		code		query		string				true	"HRM device status code (true/false)"
//	@Success	200			{object}	map[string]bool		"status: success, message: HRM device status updated"
//	@Failure	400			{object}	map[string]string	"status: error, message: Invalid request"
//	@Router		/api/v1/hrm_status [put]
func (h *HTTPHandler) SetHRMStatus(ctx *gin.Context) {

	wId, _ := parseUUID(ctx, "workout_id")

	code := ctx.Query("code")
	boolValue, boolErr := strconv.ParseBool(code)
	if boolErr != nil {
		fmt.Println(boolErr)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": boolErr.Error()})
		return
	} else {
		fmt.Println("Boolean value:", boolValue)
	}
	err := h.svc.SetHRMDevStatus(wId, boolValue)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Info("Peripheral: set hrm device status success")
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": true})

}

// SetHRMReading sets the Heart Rate Monitor (HRM) device reading from a smartwatch.
//
//	@Summary	Set HRM device reading
//	@Tags		peripheral
//	@ID			set-hrm-reading
//	@Accept		json
//	@Produce	json
//	@Param		hrm_id			path		string				true	"HRM ID"	format(uuid)
//	@Param		current_reading	query		string				true	"Current HRM reading"
//	@Success	200				{object}	map[string]string	"status: success, message: HRM reading updated"
//	@Failure	400				{object}	map[string]string	"status: error, message: Invalid request"
//	@Router		/api/v1/hrm_reading [put]
func (h *HTTPHandler) SetHRMReading(ctx *gin.Context) {
	hId, err := parseUUID(ctx, "hrm_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "error reading from smart watch"})
		return
	}
	rate := ctx.Query("current_reading")

	intValue, intErr := strconv.Atoi(rate)
	if intErr != nil {
		fmt.Println(intErr)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": intErr.Error()})
		return
	}

	if !h.svc.CheckStatusByHRMId(hId) {
		ctx.JSON(http.StatusOK, gin.H{"status": "error", "message": "hrm no such device "})
		return
	}

	err = h.svc.SetHeartRateReading(hId, intValue)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "hrm cannot read from smart watch"})
		return
	}
	log.Debug("HRM Smart Watch Reading", zap.Int("reading", intValue))
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "hrm read data from smart watch"})
}

// GetGeoStatus retrieves the Geographic (Geo) device status for a workout.
//
//	@Summary	Get Geo device status
//	@Tags		peripheral
//	@ID			get-geo-status
//	@Accept		json
//	@Produce	json
//	@Param		workout_id	path		string				true	"Workout ID"	format(uuid)
//	@Success	200			{object}	map[string]bool		"status: success, geo running: Geo device status"
//	@Failure	400			{object}	map[string]string	"status: error, message: Cannot get Geo status"
//	@Router		/api/v1/geo_status [get]
func (h *HTTPHandler) GetGeoStatus(ctx *gin.Context) {

	wId, err := parseUUID(ctx, "workout_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "cannot get geo status "})
		return
	}
	geoStatus, err := h.svc.GetGeoDevStatus(wId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "cannot get geo status"})
		return
	}
	log.Info("Peripheral: get geo device status success")
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "geo running": geoStatus})
}

// SetGeoStatus sets the Geographic (Geo) device status for a workout.
//
//	@Summary	Set Geo device status
//	@Tags		peripheral
//	@ID			set-geo-status
//	@Accept		json
//	@Produce	json
//	@Param		workout_id	path		string				true	"Workout ID"	format(uuid)
//	@Param		code		query		string				true	"Geo device status code (true/false)"
//	@Success	200			{object}	map[string]string	"status: success"
//	@Failure	400			{object}	map[string]string	"status: error, message: Invalid request"
//	@Router		/api/v1/geo_status [put]
func (h *HTTPHandler) SetGeoStatus(ctx *gin.Context) {

	wId, err := parseUUID(ctx, "workout_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "error set geo status from device"})
		return
	}
	code := ctx.Query("code")
	boolValue, boolErr := strconv.ParseBool(code)
	if boolErr != nil {
		fmt.Println(boolErr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": boolErr.Error()})
		return
	} else {
		fmt.Println("Boolean value:", boolValue)
	}
	h.svc.SetGeoDevStatus(wId, boolValue)
	log.Info("Peripheral: set geo device status success")
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

// SetGeoReading sets the Geographic (Geo) device reading for a workout.
//
//	@Summary	Set Geo device reading
//	@Tags		peripheral
//	@ID			set-geo-reading
//	@Accept		json
//	@Produce	json
//	@Param		workout_id	path		string				true	"Workout ID"	format(uuid)
//	@Param		latitude	query		string				true	"Latitude value"
//	@Param		longitude	query		string				true	"Longitude value"
//	@Success	200			{object}	map[string]string	"status: success, message: Geo reading set and location sent"
//	@Failure	400			{object}	map[string]string	"status: error, message: Invalid request"
//	@Router		/api/v1/geo_reading [put]
func (h *HTTPHandler) SetGeoReading(ctx *gin.Context) {

	latitude := ctx.Query("latitude")
	longitude := ctx.Query("longitude")
	wId, err := parseUUID(ctx, "workout_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "error read geo status from device"})
		return
	}
	var tempLastLoc LastLocation
	tempLastLoc.WorkoutID = wId
	tempLastLoc.TimeOfLocation = time.Now()

	flongitude, longFloatErr := strconv.ParseFloat(longitude, 64) // convert to float64, for float32 use '32'
	if longFloatErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "error read geo status from device"})
		return
	}

	flatitude, latFloatErr := strconv.ParseFloat(latitude, 64) // convert to float64, for float32 use '32'
	if latFloatErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "error read geo status from device"})
		return
	}

	tempLastLoc.Latitude = flatitude
	tempLastLoc.Longitude = flongitude
	err = h.svc.SetGeoLocation(wId, flongitude, flatitude)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "error read geo status from device"})
		return
	}
	log.Info("Peripheral: sending location to queue now")
	go h.svc.SendLastLocation(tempLastLoc.WorkoutID, tempLastLoc.Latitude, tempLastLoc.Longitude, tempLastLoc.TimeOfLocation)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Geo reading set and location sent"})
}

// GetGeoReading retrieves the Geographic (Geo) device reading for a workout.
//
//	@Summary	Get Geo device reading
//	@Tags		peripheral
//	@ID			get-geo-reading
//	@Accept		json
//	@Produce	json
//	@Param		workout_id	path		string				true	"Workout ID"	format(uuid)
//	@Success	200			{object}	LastLocation		"Geo device reading"
//	@Failure	400			{object}	map[string]string	"status: error, message: Failed to get Geo device"
//	@Router		/api/v1/geo_reading [get]
func (h *HTTPHandler) GetGeoReading(ctx *gin.Context) {
	wId, err := parseUUID(ctx, "workout_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to get geo device"})
		return
	}
	var tLoc LastLocation
	tLoc.TimeOfLocation, tLoc.Longitude, tLoc.Latitude, tLoc.WorkoutID, err = h.svc.GetGeoLocation(wId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to get geo device"})
		return
	}
	log.Info("Peripheral: getting location from smart watch now")
	ctx.JSON(http.StatusOK, tLoc)
}

func (h *HTTPHandler) StartBackgroundMockReading(ctx context.Context, ctx1 context.Context, wId uuid.UUID, hId uuid.UUID, longitudeStart float64, latitudeStart float64, longitudeEnd float64, latitudeEnd float64) {
	go func() {
		startLong := longitudeStart
		startLat := latitudeStart
		rand.Seed(time.Now().UnixNano())

		min := 0.0000002
		max := 0.0000005

		minHR := 80
		maxHR := 200

		// Generate a random float64 between min an
		for {
			select {
			case <-ctx1.Done(): // Check if the context is cancelled
				log.Info("Peripheral: Stopping sending data to queues...")
				return
			default:
				// Fetch the peripheral instance to check if live_data is true
				res, err := h.svc.GetLiveStatus(wId)

				if err != nil {
					h.cancelF()
				}

				if res {

					log.Info("Peripheral: Start sending data to queues...")
					if startLat <= latitudeEnd {
						randomNumber1 := randomFloat64(min, max)
						startLat += (0.0000001 + randomNumber1)
					} else {
						randomNumber1 := randomFloat64(min, max)
						startLat -= (0.0000001 + randomNumber1)
					}

					if startLong <= longitudeEnd {
						randomNumber2 := randomFloat64(min, max)
						startLong += (0.0000001 + randomNumber2)
					} else {
						randomNumber2 := randomFloat64(min, max)
						startLong -= (0.0000001 + randomNumber2)
					}

					var tLoc LastLocation
					tLoc.Longitude = startLong
					tLoc.Latitude = startLat
					tLoc.TimeOfLocation = time.Now()
					tLoc.WorkoutID = wId

					go h.svc.SendLastLocation(tLoc.WorkoutID, tLoc.Latitude, tLoc.Longitude, tLoc.TimeOfLocation)
					err2 := h.svc.SetGeoLocation(wId, tLoc.Longitude, tLoc.Latitude)
					if err2 != nil {
						log.Error("Peripheral: error sending location", zap.Error(err2))
					}

					randomInteger := rand.Intn(maxHR-minHR+1) + minHR
					currentReadingStr := fmt.Sprintf("%d", randomInteger)
					baseURL := "http://localhost:8004/api/v1/hrm_reading"
					params := url.Values{}
					params.Add("hrm_id", hId.String())
					fmt.Println(currentReadingStr)
					params.Add("current_reading", currentReadingStr)
					requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
					req, err := http.NewRequest(http.MethodPut, requestURL, bytes.NewBuffer(nil))
					if err != nil {
						// Handle error
						log.Error("Peripheral: Error creating request:", zap.Error(err))
						return
					}

					// Set headers if needed
					req.Header.Set("Content-Type", "application/json")

					// Send the request
					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						// Handle error
						log.Error("Peripheral: Error sending request:", zap.Error(err))
						return
					}
					defer resp.Body.Close()

				} else {
					log.Info("Peripheral: sending info switch is off, not publishing to queues")
					return
				}
				// Sleep for a while before printing again
				time.Sleep(1 * time.Second)
			}
		}
	}()
}
