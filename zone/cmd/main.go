package main

import (
	"github.com/CAS735-F23/macrun-teamvsl/zone/config"
	httphandler "github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/primary/http"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/primary/peripheralhandler"
	repository "github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/secondary/repository/memory"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/secondary/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/secondary/workouthandler"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/core/services"
	log "github.com/CAS735-F23/macrun-teamvsl/zone/log"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var cfg *config.AppConfiguration = config.Config

func main() {
	log.Info("Zone Manager is starting")

	// Initialize router
	router := gin.New()
	router.Use(gin.Recovery())

	// Initialize the repository
	repo := repository.NewMemoryRepository()
	repoDB := postgres.NewDBRepository(cfg.Postgres)

	// Initialize the publisher for the zone managaer to send the shelter info to the queue
	workoutDataHandler, _ := workouthandler.NewAMQPPublisher()

	// Initialize the ZoneManager service
	ZoneService, _ := services.NewZoneService(repo, repoDB, workoutDataHandler)

	// Initialize the HTTP handler with the trail manager service and the RabbitMQ handler
	ZoneManagerHTTPHandler := httphandler.NewZoneServiceHTTPHandler(router, ZoneService) // Adjusted for package

	// Set up the HTTP routes
	ZoneManagerHTTPHandler.InitRouter()
	phandler := peripheralhandler.NewAMQPHandler(ZoneService)
	phandler.InitAMQP()

	// Start the HTTP server
	err := router.Run(":" + cfg.Port)
	if err != nil {
		log.Fatal("Failed to run the server", zap.Error(err))
	}
}
