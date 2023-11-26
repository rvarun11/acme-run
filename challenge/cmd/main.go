package main

import (
	"github.com/CAS735-F23/macrun-teamvsl/challenge/config"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/docs"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/adapters/handler/amqp"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/adapters/handler/http"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/adapters/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/services"
	logger "github.com/CAS735-F23/macrun-teamvsl/challenge/log"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var cfg *config.AppConfiguration = config.Config

// @title Challenge Manager API
// @version 1.0
// @description This provides a description of API endpoints for the Challenge Manager

// @contact.name Varun Rajput
// @contact.url    https://github.com/rvarun11
// @contact.email rajpuv2@mcmaster.ca

// @query.collection.format multi
func main() {
	logger.Info("Challenge Manager is starting...")

	// Initialize router
	router := gin.Default()
	router.Use(gin.Recovery())

	// Initialize postgres repository
	store := postgres.NewRepository(cfg.Postgres)

	// Initialize challenge service
	challengeSvc := services.NewChallengeService(store)
	challengeHandler := http.NewChallengeHandler(router, challengeSvc)
	challengeHandler.InitRouter()

	// Initialize WorkoutStats Consumer
	statsConsumer := amqp.NewWorkoutStatsConsumer(cfg.RabbitMQ, challengeSvc)
	statsConsumer.InitAMQP()

	// Initialize badge service
	// badgeSvc := services.NewBadgeService(store)
	// badgeHandler := http.NewBadgeHandler(router, *badgeSvc)
	// badgeHandler.InitRouter()

	// Swagger Support
	docs.SwaggerInfo.Host = "localhost:" + cfg.Port
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start Server
	router.Run(":" + cfg.Port)
}
