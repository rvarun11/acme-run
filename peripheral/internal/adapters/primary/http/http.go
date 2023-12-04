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

	"github.com/CAS735-F23/macrun-teamvsl/peripheral/config"
	rabbitmqhandler "github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/adapters/secondary/amqp"
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

	// HRM
	router.POST("/peripheral/hrm", handler.connectHRM)
	router.PUT("/peripheral/hrm/:hrm_id", handler.disconnectHRM)
	router.GET("/peripheral/hrm", handler.getHRMReading)

	router.PUT("/hrm/:hrm_id", handler.SetHRMReading)
	router.PUT("/geo/:geo_id", handler.SetGeoReading)

}

// ConnectHRM connects to a Heart Rate Monitor (HRM) device.
//
//	@Summary	Connect to HRM device
//	@Tags		peripheral
//	@ID			connect-hrm
//	@Accept		json
//	@Produce	json
//	@Param		connectData	body	BindPeripheralData	true	"Connect Peripheral Data"
//	@Success	200			"connect to hrm success: true"
//	@Failure	400			"Bad Request with error details"
//	@Failure	500			"Internal error with error details"
//	@Router		/api/v1/peripheral/hrm [post]
func (h *HTTPHandler) connectHRM(ctx *gin.Context) {
	var req BindPeripheralData
	var err error
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request parameters",
		})
		return
	}

	if !h.svc.CheckStatusByHRMId(req.HRMId) {
		err = h.svc.CreatePeripheral(req.PlayerID, req.HRMId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "could not connect HRM, something went wrong",
			})
			return
		}
	}

	err2 := h.svc.SetHRMDevStatusByHRMId(req.HRMId, true)
	if err2 != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not connect HRM, something went wron",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "hrm connected successfully"})
}

// DisconnectHRM disconnects a Heart Rate Monitor (HRM) device.
//
//	@Summary	Disconnect HRM device
//	@Tags		peripheral
//	@ID			disconnect-hrm
//	@Accept		json
//	@Produce	json
//	@Param		disconnectData	body	BindPeripheralData	true	"Disconnect Peripheral Data"
//	@Success	200				"disconnected HRM"
//	@Failure	400				"error invalid request"
//	@Failure	500				"cannot disconnect HRM"
//	@Router		/api/v1/peripheral/hrm/{hrm_id} [put]
func (h *HTTPHandler) disconnectHRM(ctx *gin.Context) {

	hrmIdStr := ctx.Param("hrm_id")
	hrmId, err := uuid.Parse(hrmIdStr)
	if err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to disconnect to hrm, invalid hrm id"})
		return

	}
	hrmStatus, err := h.svc.GetHRMDevStatusByHRMId(hrmId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to disconnect to hrm, bad hrm request"})
		return
	}
	if !hrmStatus {
		ctx.JSON(http.StatusOK, gin.H{"message": "hrm already disconnected"})
		return
	}

	err2 := h.svc.SetHRMDevStatusByHRMId(hrmId, false)
	if err2 != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot disconnect hrm, something went wrong",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "hrm disconnected successfully"})

}

// BindPeripheralToData connects to an HRM device and binds it to a workout
//
//	@Summary	Bind peripheral to a workout
//	@Tags		peripheral
//	@ID			bind-peripheral
//	@Accept		json
//	@Produce	json
//	@Param		bindData	body	BindPeripheralData	true	"Bind Peripheral Data"
//	@Success	200			"binding workout done"
//	@Failure	400			"error unbind with message"
//	@Router		/api/v1/peripheral [post]
func (h *HTTPHandler) BindPeripheralToData(ctx *gin.Context) {

	var bindDataInstance BindPeripheralData
	if err := ctx.ShouldBindJSON(&bindDataInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot bind peripheral data, something went wrong"})
		return
	}

	err := h.svc.BindPeripheral(bindDataInstance.PlayerID, bindDataInstance.WorkoutID, bindDataInstance.HRMId, bindDataInstance.HRMConnect, bindDataInstance.SendLiveLocationToTrailManager)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to bind workout",
		})
		return
	} else {

		h.hLiveCount += 1
		h.bCtx, h.cancelF = context.WithCancel(context.Background())

		longitudeStart, latitudeStart, longitudeEnd, latitudeEnd, err := h.svc.GetTrailLocationInfo(bindDataInstance.TrailOfWorkout)
		if err != nil {
			log.Error("peripheral: failed to get trail location, using default info now", zap.Error(err))
			longitudeStart = -79.919390
			latitudeStart = 43.257715
			longitudeEnd = 43.258012
			latitudeEnd = -79.910866
		}

		h.svc.SetLiveStatus(bindDataInstance.WorkoutID, true)
		h.StartBackgroundMockReading(ctx, h.bCtx, bindDataInstance.WorkoutID, bindDataInstance.HRMId, longitudeStart, latitudeStart, longitudeEnd, latitudeEnd)
		log.Info("peripheral: bound to workout successfully", zap.Any("workout_id", bindDataInstance.WorkoutID), zap.Any("hrm_id", bindDataInstance.HRMId))
		ctx.JSON(http.StatusOK, gin.H{
			"message": "binding workout successful",
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
//	@Param		unbindData	body	UnbindPeripheralData	true	"Unbind Peripheral Data"
//	@Success	200			"success message: Unbind the data"
//	@Failure	400			"error message: invalid request with details"
//	@Failure	500			"error message: with details"
//	@Router		/api/v1/peripheral [put]
func (h *HTTPHandler) UnbindPeripheralToData(ctx *gin.Context) {
	var req UnbindPeripheralData
	var err error
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request parameters",
		})
		return
	}

	err = h.svc.SetLiveStatus(req.WorkoutID, false)
	if err != nil {
		log.Error("peripheral: failed to set live status of publising ", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to unbind workout",
		})
		return
	}

	h.hLiveCount -= 1
	if h.hLiveCount == 0 {
		h.cancelF()
	}
	log.Info("peripheral: bound to workout successfully", zap.Any("workout_id", req.WorkoutID))
	ctx.JSON(http.StatusOK, gin.H{
		"message": "peripheral unbound from workout"})
}

// getHRMReading retrieves Heart Rate Monitor (HRM) reading data.
//
//	@Summary	Get average heart rate
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout_id"})
		return
	}

	hrType := ctx.Query("type")
	if hrType == "avg" {
		// TODO: This should be returning as per workout
		var tLoc LastHR
		var err error
		tLoc.HRMID, tLoc.TimeOfLocation, tLoc.HeartRate, err = h.svc.GetHRMAvgReading(wId)
		if err != nil {
			log.Error("peripheral: failed to read from device failure ", zap.Error(err))
			ctx.JSON(http.StatusOK, gin.H{
				"message": "heart rate record not found",
			})
			return
		}
		avgRate := AverageHeartRate{}
		avgRate.WorkoutID = wId
		avgRate.AverageHeartRate = uint8(tLoc.HeartRate)

		jsonData, err := json.Marshal(avgRate)
		if err != nil {
			log.Error("peripheral: failed to read from device,failed to marshal json ", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not marshal JSON for average heart rate: " + err.Error(),
			})
			return
		}
		log.Info("average hrm read successfully", zap.Any("value", tLoc.HeartRate))
		ctx.Writer.Header().Set("Content-Type", "application/json")
		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Writer.Write(jsonData)
	} else if hrType == "normal" {
		var tLoc LastHR
		var err error
		tLoc.HRMID, tLoc.TimeOfLocation, tLoc.HeartRate, err = h.svc.GetHRMReading(wId)
		if err != nil {
			log.Error("peripheral: failed to read from device failure ", zap.Error(err))
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "reading from hrm failed",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"reading": tLoc})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "reading from device failure, invalid argument of type",
		})
	}
}

func (h *HTTPHandler) GetHRMStatus(ctx *gin.Context) {

	wId, err := parseUUID(ctx, "workout_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout_id"})
		return
	}

	tStatus, err := h.svc.GetHRMDevStatus(wId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to get hrm status"})
		return
	}
	log.Info("peripheral: get geo device status success")
	ctx.JSON(http.StatusOK, gin.H{"message": tStatus})

}

func (h *HTTPHandler) SetHRMStatus(ctx *gin.Context) {

	wId, _ := parseUUID(ctx, "workout_id")

	code := ctx.Query("code")
	boolValue, boolErr := strconv.ParseBool(code)
	if boolErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": boolErr.Error()})
		return
	}
	err := h.svc.SetHRMDevStatus(wId, boolValue)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Info("peripheral: set hrm device status success")
	ctx.JSON(http.StatusOK, gin.H{"message": true})

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
//	@Router		/api/v1/hrm/:hrm_id [put]
func (h *HTTPHandler) SetHRMReading(ctx *gin.Context) {

	hrmIdStr := ctx.Param("hrm_id")
	hId, err := uuid.Parse(hrmIdStr)
	if err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to disconnect to hrm, invalid hrm id"})
		return

	}

	rate := ctx.Query("current_reading")
	intValue, intErr := strconv.Atoi(rate)
	if intErr != nil {
		fmt.Println(intErr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": intErr.Error()})
		return
	}

	if !h.svc.CheckStatusByHRMId(hId) {
		ctx.JSON(http.StatusOK, gin.H{"error": "hrm no such device "})
		return
	}

	err = h.svc.SetHeartRateReading(hId, intValue)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "hrm cannot read from smart watch"})
		return
	}
	log.Debug("HRM Smart Watch Reading", zap.Int("reading", intValue))
	ctx.JSON(http.StatusOK, gin.H{"message": "hrm read data from smart watch and set "})
}

func (h *HTTPHandler) GetGeoStatus(ctx *gin.Context) {

	wId, err := parseUUID(ctx, "workout_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot get geo status "})
		return
	}
	geoStatus, err := h.svc.GetGeoDevStatus(wId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot get geo status"})
		return
	}
	log.Info("peripheral: get geo device status success")
	ctx.JSON(http.StatusOK, gin.H{"message": "geo status is " + strconv.FormatBool(geoStatus)})
}

func (h *HTTPHandler) SetGeoStatus(ctx *gin.Context) {

	wId, err := parseUUID(ctx, "workout_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error set geo status from device"})
		return
	}
	code := ctx.Query("code")
	boolValue, boolErr := strconv.ParseBool(code)
	if boolErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": boolErr.Error()})
		return
	}
	h.svc.SetGeoDevStatus(wId, boolValue)
	log.Debug("peripheral: set geo device status success")
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

// SetGeoReading sets the Geographic (Geo) device reading for a workout.
//
//	@Summary	Set live location (geo reading)
//	@Tags		peripheral
//	@ID			set-geo-reading
//	@Accept		json
//	@Produce	json
//	@Param		geo_id		path	string	true	"Workout ID"	format(uuid)
//	@Param		latitude	query	string	true	"Latitude value"
//	@Param		longitude	query	string	true	"Longitude value"
//	@Success	200			" message, geo reading set and location sent"
//	@Failure	400			"error message with details"
//	@Router		/api/v1/geo/:geo_id [put]
func (h *HTTPHandler) SetGeoReading(ctx *gin.Context) {

	workoutIdStr := ctx.Param("workout_id")
	wId, err := uuid.Parse(workoutIdStr)
	if err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to set geo reading, bad workout id"})
		return

	}
	latitude := ctx.Query("latitude")
	longitude := ctx.Query("longitude")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error set geo status from device"})
		return
	}
	var tempLastLoc LastLocation
	tempLastLoc.WorkoutID = wId
	tempLastLoc.TimeOfLocation = time.Now()

	flongitude, longFloatErr := strconv.ParseFloat(longitude, 64) // convert to float64, for float32 use '32'
	if longFloatErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error set geo status from device"})
		return
	}

	flatitude, latFloatErr := strconv.ParseFloat(latitude, 64) // convert to float64, for float32 use '32'
	if latFloatErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error  set geo status from device"})
		return
	}

	tempLastLoc.Latitude = flatitude
	tempLastLoc.Longitude = flongitude
	err = h.svc.SetGeoLocation(wId, flongitude, flatitude)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error set geo status from device"})
		return
	}
	log.Info("peripheral: sending location to queue now")
	go h.svc.SendLastLocation(tempLastLoc.WorkoutID, tempLastLoc.Latitude, tempLastLoc.Longitude, tempLastLoc.TimeOfLocation)
	ctx.JSON(http.StatusOK, gin.H{"message": "geo reading set and location sent"})
}

func (h *HTTPHandler) GetGeoReading(ctx *gin.Context) {
	wId, err := parseUUID(ctx, "workout_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to get geo device"})
		return
	}
	var tLoc LastLocation
	tLoc.TimeOfLocation, tLoc.Longitude, tLoc.Latitude, tLoc.WorkoutID, err = h.svc.GetGeoLocation(wId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to get geo device"})
		return
	}
	log.Info("peripheral: getting location from smart watch now")
	ctx.JSON(http.StatusOK, tLoc)
}

/*
NOTE: StartBackgroundMockReading for mocking/simulating external fitness devices:
It uses the APIs provided by the Peripheral Service to generate heartrate & geo data
*/
var port string = config.Config.Port

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
				log.Debug("peripheral: stopping sending data to queues...")
				return
			default:
				// Fetch the peripheral instance to check if live_data is true
				res, err := h.svc.GetLiveStatus(wId)

				if err != nil {
					h.cancelF()
				}

				if res {

					log.Debug("Peripheral: Start sending data to queues...")
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
					baseURL := "http://localhost:" + port + "/api/v1/hrm/" + hId.String()
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
					log.Debug("Peripheral: sending info switch is off, not publishing to queues")
					return
				}
				// Sleep for a while before printing again
				time.Sleep(1 * time.Second)
			}
		}
	}()
}
