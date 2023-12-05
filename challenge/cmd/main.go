package main

import (
	"github.com/CAS735-F23/macrun-teamvsl/challenge/config"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/docs"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/adapters/primary/amqp"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/adapters/primary/http"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/adapters/secondary/repository/postgres"
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
	logger.Info("challenge manager is starting...")

	// Initialize router
	router := gin.New()
	router.Use(gin.Recovery())

	// Initialize postgres repository
	store := postgres.NewRepository(cfg.Postgres)

	// Initialize challenge service
	challengeSvc := services.NewChallengeService(store)
	challengeHandler := http.NewChallengeHandler(router, challengeSvc)
	challengeHandler.InitRouter()

	// Initialize workout stats consumer
	workoutStatsConsumer := amqp.NewWorkoutStatsConsumer(cfg.RabbitMQ, challengeSvc)
	workoutStatsConsumer.InitAMQP()

	// Swagger support
	docs.SwaggerInfo.Host = "localhost:" + cfg.Port
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start Server
	router.Run(":" + cfg.Port)
}
