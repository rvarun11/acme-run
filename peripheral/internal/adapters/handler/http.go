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
	router.POST("/peripheral", handler.CreatePeripheralDevice)
	router.POST("/peripheral_bind", handler.BindPeripheralToData)
	router.POST("/peripheral_unbind", handler.UnbindPeripheralToData)
	router.PUT("/hrm_status", handler.SetHRMStatus)
	router.POST("/hrm_connect", handler.ConnectHRM)
	router.GET("/hrm_status", handler.GetHRMStatus)
	router.GET("/hrm_reading", handler.GetHRMReading)
	router.PUT("/hrm_reading", handler.SetHRMReading)
	router.PUT("/geo_status", handler.SetGeoStatus)
	router.GET("/geo_status", handler.GetGeoStatus)
	router.GET("/geo_reading", handler.GetGeoReading)
	router.PUT("/geo_reading", handler.SetGeoReading)

}

func (h *HTTPHandler) ConnectHRM(ctx *gin.Context) {

	hId, err1 := parseUUID(ctx, "hrm_id")
	code := ctx.Query("code")
	if err1 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"failed to connect to hrm device failed, error 1": err1,
		})
		return
	}

	boolValue, boolErr := strconv.ParseBool(code)
	if boolErr != nil {
		fmt.Println(boolErr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": boolErr.Error()})
		return
	} else {
		fmt.Println("Boolean value:", boolValue)
	}

	if h.svc.CheckStatusByHRMId(hId) {

	} else {
		pId, err1 := parseUUID(ctx, "player_id")
		if err1 != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"failed to connect to hrm device failed, error 1": err1,
			})
			return
		}
		h.svc.CreatePeripheral(pId, hId)
	}

	h.svc.SetHRMDevStatusByHRMId(hId, boolValue)
	if boolValue == true {
		ctx.JSON(http.StatusOK, gin.H{"connect to hrm success": true})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"disconnect to hrm success": true})
	}
}

func (h *HTTPHandler) CreatePeripheralDevice(ctx *gin.Context) {
	// pId, err1 := parseUUID(ctx, "player_id")
	// hId, err2 := parseUUID(ctx, "hrm_id")
	// if err2 != nil || err1 != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"creating a new peripheral device failed, error 1": err1, "error 2": err2,
	// 	})
	// 	return
	// }

	var cDataInstance BindPeripheralData
	if err := ctx.ShouldBindJSON(&cDataInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.svc.CreatePeripheral(cDataInstance.WorkoutID, cDataInstance.HRMId)
	ctx.JSON(http.StatusOK, gin.H{"device creation": true})
}

func (h *HTTPHandler) BindPeripheralToData(ctx *gin.Context) {

	// pId, err := parseUUID(ctx, "player_id")
	// if err != nil {

	// }
	// wId, err := parseUUID(ctx, "workout_id")
	// if err != nil {

	// }
	// hId, err := parseUUID(ctx, "hrm_id")
	// if err != nil {

	// }
	// tConnected := ctx.Query("hrm_connected")
	// connected, boolErr := strconv.ParseBool(tConnected)
	// if boolErr != nil {
	// 	fmt.Println(boolErr)
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": boolErr.Error()})
	// 	return
	// } else {
	// 	fmt.Println("Boolean value:", connected)
	// }

	// tBroadcast := ctx.Query("send_live_location_to_trail_manager")
	// broadcast, boolErr := strconv.ParseBool(tBroadcast)
	// if boolErr != nil {
	// 	fmt.Println(boolErr)
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": boolErr.Error()})
	// 	return
	// } else {
	// 	fmt.Println("Boolean value:", broadcast)
	// }

	var bindDataInstance BindPeripheralData
	if err := ctx.ShouldBindJSON(&bindDataInstance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.svc.BindPeripheral(bindDataInstance.PlayerID, bindDataInstance.WorkoutID, bindDataInstance.HRMId, bindDataInstance.HRMConnected, bindDataInstance.SendLiveLocationToTrailManager)
	// h.svc.BindPeripheral(pId, wId, hId, connected, broadcast)
	h.hLiveCount += 1

	// if connected {
	// go h.rabbitMQHandler.SendLastHR(wId)
	ctx.JSON(http.StatusBadRequest, gin.H{
		"binding successful": true,
	})
	// } else {
	// 	h.svc.DisconnectPeripheral(bindDataInstance.WorkoutID)
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"binding successful": false,
	// 	})
	// }

	// Create a context that can be cancelled
	h.bCtx, h.cancelF = context.WithCancel(context.Background())

	longitudeStart := 40.0
	latitudeStart := 40.0
	longitudeEnd := 50.0
	latitudeEnd := 50.0
	// liveDataSw := true
	h.svc.SetLiveSw(bindDataInstance.WorkoutID, true)
	// Start the background printing
	h.StartBackgroundMockTesting(ctx, h.bCtx, bindDataInstance.WorkoutID, bindDataInstance.HRMId, longitudeStart, latitudeStart, longitudeEnd, latitudeEnd)

}

func (h *HTTPHandler) UnbindPeripheralToData(ctx *gin.Context) {

	wId, err := parseUUID(ctx, "workout_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"unbind": false})
		return
	}
	h.svc.SetLiveSw(wId, false)
	h.hLiveCount -= 1
	if h.hLiveCount == 0 {
		h.cancelF()
	}
	h.svc.SetHRMDevStatus(wId, false)
	h.svc.SetGeoDevStatus(wId, false)
	h.svc.DisconnectPeripheral(wId)
	ctx.JSON(http.StatusOK, gin.H{
		"unbind": true})
}

func (h *HTTPHandler) GetHRMReading(ctx *gin.Context) {
	hId, err := parseUUID(ctx, "hrm_id")
	// wId := ctx.Query("workout_id")
	if err != nil {

	}
	ReadingType := ctx.Query("type")
	if ReadingType == "avg" {
		ctx.JSON(http.StatusOK, gin.H{"reading": h.svc.GetHRMAvgReading(hId)})
	} else if ReadingType == "normal" {
		ctx.JSON(http.StatusOK, h.svc.GetHRMReading(hId))
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"reading from device success": false})
	}
}

func (h *HTTPHandler) GetHRMStatus(ctx *gin.Context) {

	// wId, err := parseUUID(ctx, "workout_id")
	wId, err := parseUUID(ctx, "hrm_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"fail to get hrm status": false})
	}
	ctx.JSON(http.StatusOK, gin.H{"hrm status value": h.svc.GetHRMDevStatusByHRMId(wId)})

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
	h.svc.SetHRMDevStatus(wId, boolValue)
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
	} else {
		log.Debug("HRM Smart Watch Reading", zap.Int("reading", intValue))
	}
	h.svc.SetHeartRateReading(hId, intValue)
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *HTTPHandler) GetGeoStatus(ctx *gin.Context) {

	wId, err := parseUUID(ctx, "workout_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"cannot get geo status ": false})
		return
	}
	ctx.JSON(http.StatusOK, h.svc.GetGeoDevStatus(wId))
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error read geo status from device, success": false})
	}
	var tempLastLoc LastLocation
	tempLastLoc.WorkoutID = wId
	tempLastLoc.TimeOfLocation = time.Now()

	flongitude, longFloatErr := strconv.ParseFloat(longitude, 64) // convert to float64, for float32 use '32'
	if longFloatErr != nil {
		fmt.Println(longFloatErr)
	} else {
		// fmt.Println("Float value:", flongitude)
	}

	flatitude, latFloatErr := strconv.ParseFloat(latitude, 64) // convert to float64, for float32 use '32'
	if latFloatErr != nil {
		fmt.Println(latFloatErr)
	} else {
		// fmt.Println("Float value:", flatitude)
	}
	tempLastLoc.Latitude = flatitude
	tempLastLoc.Longitude = flongitude
	h.svc.SetGeoLocation(wId, flongitude, flatitude)
	// Now trigger the RabbitMQHandler to send the updated location
	log.Info("sending location to queue now")
	go h.rabbitMQHandler.SendLastLocation(tempLastLoc)
	log.Info("sent location to queue now")
	fmt.Println(h.svc.GetGeoLocation(wId))
	ctx.JSON(http.StatusOK, gin.H{"message": "Geo reading set and location sent"})

	// ctx.JSON(http.StatusOK, h.svc.SetGeoLocation(wId, longitude, latitude))
}

func (h *HTTPHandler) GetGeoReading(ctx *gin.Context) {
	wId, err := parseUUID(ctx, "workout_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error read geo location, success": false})
	}
	ctx.JSON(http.StatusOK, h.svc.GetGeoLocation(wId))
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
				if h.svc.GetLiveSw(wId) {

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
