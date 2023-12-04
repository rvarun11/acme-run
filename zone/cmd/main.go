package main

import (
	"github.com/CAS735-F23/macrun-teamvsl/zone/config"
	"github.com/CAS735-F23/macrun-teamvsl/zone/docs"
	amqpPrimary "github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/primary/amqp"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/primary/http"
	amqpSecondary "github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/secondary/amqp"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/secondary/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/core/services"
	logger "github.com/CAS735-F23/macrun-teamvsl/zone/log"
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
	logger.Info("zone manager is starting...")

	// Initialize router
	router := gin.New()
	router.Use(gin.Recovery())

	// Initialize the repository
	// repo := repository.NewMemoryRepository()
	repo := postgres.NewRepository(cfg.Postgres)

	// Initialize shelter distance publisher
	shelterDistancePublisher := amqpSecondary.NewShelterDistancePublisher(cfg.RabbitMQ)

	// Initialize the zone manager
	zoneSvc, _ := services.NewZoneService(repo, shelterDistancePublisher)
	zoneHandler := http.NewZoneHandler(router, zoneSvc)
	zoneHandler.InitRouter()

	// Initialize location consumer
	locationConsumer := amqpPrimary.NewLocationConsumer(cfg.RabbitMQ, zoneSvc)
	locationConsumer.InitAMQP()

	// Swagger support
	docs.SwaggerInfo.Host = "localhost:" + cfg.Port
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the HTTP server
	router.Run(":" + cfg.Port)
}
