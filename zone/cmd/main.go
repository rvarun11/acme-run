package main

import (
	"github.com/CAS735-F23/macrun-teamvsl/zone/config"
	"github.com/CAS735-F23/macrun-teamvsl/zone/docs"
	httphandler "github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/primary/http"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/primary/peripheralhandler"
	repository "github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/secondary/repository/memory"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/secondary/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/secondary/workouthandler"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/core/services"
	log "github.com/CAS735-F23/macrun-teamvsl/zone/log"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var cfg *config.AppConfiguration = config.Config

// @title Zone Service API
// @version 1.0
// @description This provides a description of API endpoints for the zone

// @contact.name Liuyin Shi
// @contact.url    https://github.com/XIAOKAOBO
// @contact.email shil9@mcmaster.ca

// @query.collection.format multi

func main() {
	log.Info("zone manager is starting...")

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

	docs.SwaggerInfo.Host = "localhost:" + cfg.Port
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the HTTP server
	router.Run(":" + cfg.Port)
}
