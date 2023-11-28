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

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Binding workout done",
		})
		h.hLiveCount += 1
		h.bCtx, h.cancelF = context.WithCancel(context.Background())

		longitudeStart, latitudeStart, longitudeEnd, latitudeEnd, errL := h.svc.GetTrailLocationInfo(bindDataInstance.TrailOfWorkout)
		if errL != nil {
			log.Error("failed to get trail location, using default info now", zap.Error(errL))
			longitudeStart = -79.919390
			latitudeStart = 43.257715
			longitudeEnd = 43.258012
			latitudeEnd = -79.910866
		}

		h.svc.SetLiveStatus(bindDataInstance.WorkoutID, true)
		h.StartBackgroundMockReading(ctx, h.bCtx, bindDataInstance.WorkoutID, bindDataInstance.HRMId, longitudeStart, latitudeStart, longitudeEnd, latitudeEnd)
	}
}

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

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Unbind the data"})
}

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
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "error", "message": "reading from device failure",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"reading": tLoc})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error", "message": "reading from device failure",
		})
	}
}

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
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": tStatus})

}

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
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": true})

}

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
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "geo running": geoStatus})
}

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
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

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
	log.Info("sending location to queue now")
	go h.svc.SendLastLocation(tempLastLoc.WorkoutID, tempLastLoc.Latitude, tempLastLoc.Longitude, tempLastLoc.TimeOfLocation)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Geo reading set and location sent"})
}

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
				fmt.Println("Stopping background printing...")
				return
			default:
				// Fetch the peripheral instance to check if live_data is true
				res, err := h.svc.GetLiveStatus(wId)

				if err != nil {
					h.cancelF()
				}

				if res {

					fmt.Println("hello")
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
						log.Error("error in sending location", zap.Error(err2))
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
						fmt.Println("Error creating request:", err)
						return
					}

					// Set headers if needed
					req.Header.Set("Content-Type", "application/json")

					// Send the request
					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						// Handle error
						fmt.Println("Error sending request:", err)
						return
					}
					defer resp.Body.Close()

				} else {
					fmt.Println("LiveData is false, stopping background printing...")
					return
				}
				// Sleep for a while before printing again
				time.Sleep(1 * time.Second)
			}
		}
	}()
}
