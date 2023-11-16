package handler

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

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

// LS-TODO: An adapter cannot talk to another adapter, rabbitMQ handler cannot be here
type HTTPHandler struct {
	gin             *gin.Engine
	svc             *services.PeripheralService
	rabbitMQHandler *RabbitMQHandler
	cancelF         context.CancelFunc
	bCtx            context.Context
	hLiveCount      int
}

func NewPeripheralServiceHTTPHandler(gin *gin.Engine, PeripheralService *services.PeripheralService, rabbitMQHandler *RabbitMQHandler) *HTTPHandler {
	return &HTTPHandler{
		gin:             gin,
		svc:             PeripheralService,
		rabbitMQHandler: rabbitMQHandler,
	}
}

func (handler *HTTPHandler) InitRouter() {

	router := handler.gin.Group("/api/v1")

	// VR TODO: Fix API endpoints
	router.POST("/peripheral", handler.CreatePeripheralDevice)
	router.POST("/peripheral_bind", handler.BindPeripheralToData)
	router.POST("/peripheral_unbind", handler.UnbindPeripheralToData)
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

// LS-TODO: Remove bool connect
// LS-TODO: Update all the handler functions to be private, for example:
func (h *HTTPHandler) connectHRM(ctx *gin.Context) {
	var req BindPeripheralData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	if !h.svc.CheckStatusByHRMId(req.HRMId) {
		err1 := h.svc.CreatePeripheral(req.PlayerID, req.HRMId)
		if err1 != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "cannot connect to hrm",
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
			"error": "invalid request",
		})
		return
	}

	err2 := h.svc.SetHRMDevStatusByHRMId(req.HRMId, false)
	if err2 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot disconnect to hrm",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"disconnect to hrm success": true})
}

func (h *HTTPHandler) CreatePeripheralDevice(ctx *gin.Context) {

	var cDataInstance BindPeripheralData
	if err := ctx.ShouldBindJSON(&cDataInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.svc.CreatePeripheral(cDataInstance.WorkoutID, cDataInstance.HRMId)
	ctx.JSON(http.StatusOK, gin.H{"device creation": true})
}

func (h *HTTPHandler) BindPeripheralToData(ctx *gin.Context) {

	var bindDataInstance BindPeripheralData
	if err := ctx.ShouldBindJSON(&bindDataInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.BindPeripheral(bindDataInstance.PlayerID, bindDataInstance.WorkoutID, bindDataInstance.HRMId, bindDataInstance.HRMConnected, bindDataInstance.SendLiveLocationToTrailManager)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to Bind workout",
		})
		return
	} else {

		ctx.JSON(http.StatusOK, gin.H{
			"success": "Binding workout done",
		})
		h.hLiveCount += 1
		h.bCtx, h.cancelF = context.WithCancel(context.Background())
		longitudeStart := 40.0
		latitudeStart := 40.0
		longitudeEnd := 50.0
		latitudeEnd := 50.0
		h.svc.SetLiveSw(bindDataInstance.WorkoutID, true)
		h.StartBackgroundMockTesting(ctx, h.bCtx, bindDataInstance.WorkoutID, bindDataInstance.HRMId, longitudeStart, latitudeStart, longitudeEnd, latitudeEnd)
	}
}

func (h *HTTPHandler) UnbindPeripheralToData(ctx *gin.Context) {
	var req UnbindPeripheralData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	err1 := h.svc.SetLiveSw(req.WorkoutID, false)
	if err1 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to unbind",
		})
		return
	}

	h.hLiveCount -= 1
	if h.hLiveCount == 0 {
		h.cancelF()
	}

	err1 = h.svc.SetHRMDevStatus(req.WorkoutID, false)
	if err1 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to unbind",
		})
		return
	}

	err1 = h.svc.SetGeoDevStatus(req.WorkoutID, false)
	if err1 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to unbind",
		})
		return
	}

	err1 = h.svc.DisconnectPeripheral(req.WorkoutID)
	if err1 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to unbind",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": "Unbind the data"})
}

func (h *HTTPHandler) getHRMReading(ctx *gin.Context) {
	wId, err1 := parseUUID(ctx, "workout_id")
	if err1 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "get hr invalid request"})
	}

	hrType := ctx.Query("type")
	if hrType == "avg" {
		// TODO: This should be returning as per workout
		tLoc, err := h.svc.GetHRMAvgReading(wId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "reading from device failure",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"reading": tLoc})
	} else if hrType == "normal" {
		tLoc, err := h.svc.GetHRMAvgReading(wId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "reading from device failure",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"reading": tLoc})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "reading from device failure",
		})
	}
}

func (h *HTTPHandler) GetHRMStatus(ctx *gin.Context) {

	wId, err := parseUUID(ctx, "workout_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Fail to get hrm status"})
		return
	}

	tStatus, err := h.svc.GetHRMDevStatus(wId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Fail to get hrm status"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"hrm running": tStatus})

}

func (h *HTTPHandler) SetHRMStatus(ctx *gin.Context) {

	wId, err := parseUUID(ctx, "workout_id")
	if err != nil {

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
	err = h.svc.SetHRMDevStatus(wId, boolValue)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true})

}

func (h *HTTPHandler) SetHRMReading(ctx *gin.Context) {
	hId, err := parseUUID(ctx, "hrm_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error reading from smart watch success": false})
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
	ctx.JSON(http.StatusOK, gin.H{"success": "hrm read data from smart watch"})
}

func (h *HTTPHandler) GetGeoStatus(ctx *gin.Context) {

	wId, err := parseUUID(ctx, "workout_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot get geo status "})
		return
	}
	geoStatus, err := h.svc.GetGeoDevStatus(wId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot get geo status "})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"geo running": geoStatus})
}

func (h *HTTPHandler) SetGeoStatus(ctx *gin.Context) {

	wId, err := parseUUID(ctx, "workout_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error set geo status from device, success": false})
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
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *HTTPHandler) SetGeoReading(ctx *gin.Context) {

	latitude := ctx.Query("latitude")
	longitude := ctx.Query("longitude")
	wId, err := parseUUID(ctx, "workout_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error read geo status from device"})
		return
	}
	var tempLastLoc LastLocation
	tempLastLoc.WorkoutID = wId
	tempLastLoc.TimeOfLocation = time.Now()

	flongitude, longFloatErr := strconv.ParseFloat(longitude, 64) // convert to float64, for float32 use '32'
	if longFloatErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error read geo status from device"})
		return
	}

	flatitude, latFloatErr := strconv.ParseFloat(latitude, 64) // convert to float64, for float32 use '32'
	if latFloatErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error read geo status from device"})
		return
	}

	tempLastLoc.Latitude = flatitude
	tempLastLoc.Longitude = flongitude
	err = h.svc.SetGeoLocation(wId, flongitude, flatitude)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error read geo status from device"})
		return
	}
	log.Info("sending location to queue now")
	go h.rabbitMQHandler.SendLastLocation(tempLastLoc)
	ctx.JSON(http.StatusOK, gin.H{"message": "Geo reading set and location sent"})
}

func (h *HTTPHandler) GetGeoReading(ctx *gin.Context) {
	wId, err := parseUUID(ctx, "workout_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get geo device"})
		return
	}
	tLoc, err := h.svc.GetGeoLocation(wId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get geo device"})
		return
	}
	ctx.JSON(http.StatusOK, tLoc)
}

func (h *HTTPHandler) StartBackgroundMockTesting(ctx context.Context, ctx1 context.Context, wId uuid.UUID, hId uuid.UUID, longitudeStart float64, latitudeStart float64, longitudeEnd float64, latitudeEnd float64) {
	go func() {
		startLong := longitudeStart
		startLat := latitudeStart
		rand.Seed(time.Now().UnixNano())

		min := 0.01
		max := 0.05

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
				res, err := h.svc.GetLiveSw(wId)

				if err != nil {
					h.cancelF()
				}

				if res {

					fmt.Println("hello")
					if startLat <= latitudeEnd {
						randomNumber1 := randomFloat64(min, max)
						startLat += (0.05 + randomNumber1)
					}

					if startLong <= longitudeEnd {
						randomNumber2 := randomFloat64(min, max)
						startLong += (0.05 + randomNumber2)
					}

					var tLoc LastLocation
					tLoc.Longitude = startLong
					tLoc.Latitude = startLat
					tLoc.TimeOfLocation = time.Now()
					tLoc.WorkoutID = wId

					go h.rabbitMQHandler.SendLastLocation(tLoc)
					err2 := h.svc.SetGeoLocation(wId, tLoc.Longitude, tLoc.Latitude)
					if err2 != nil {

					}

					randomInteger := rand.Intn(maxHR-minHR+1) + minHR
					currentReadingStr := fmt.Sprintf("%d", randomInteger)
					baseURL := "http://localhost:8004/api/v1/hrm_reading"
					params := url.Values{}
					params.Add("hrm_id", hId.String())
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
