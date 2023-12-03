package main

import (
	"fmt"
	"os"

	"github.com/CAS735-F23/macrun-teamvsl/peripheral/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/CAS735-F23/macrun-teamvsl/peripheral/config"
	httphandler "github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/adapters/handler/primary/http"
	rabbitmqhandler "github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/adapters/handler/primary/rabbitmq"
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/adapters/repository"
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/adapters/secondary/clients"
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/services"
	log "github.com/CAS735-F23/macrun-teamvsl/peripheral/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var cfg *config.AppConfiguration = config.Config

// @title Peripheral Service API
// @version 1.0
// @description This provides a description of API endpoints for the Peripheral Service

// @contact.name Liuyin Shi
// @contact.url    https://github.com/XIAOKAOBO
// @contact.email shil9>@mcmaster.ca

// @query.collection.format multi
func main() {
	log.Info("Peripheral Service is starting")

	// Initialize router
	router := gin.New()
	router.Use(gin.Recovery())

	// Initialize the repository
	repo := repository.NewMemoryRepository()

	amqpURL := "amqp://" + cfg.RabbitMQ.User + ":" +
		cfg.RabbitMQ.Password + "@" + cfg.RabbitMQ.Host + ":" + cfg.RabbitMQ.Port + "/"
	// Initialize the RabbitMQ handler with the Peripheral service and the AMQP URL

	peripheralAMQPHandler, err1 := rabbitmqhandler.NewRabbitMQHandler(amqpURL) // Adjusted for package

	client := clients.NewZoneServiceClient()

	// Initialize the Peripheral service
	peripheralService := services.NewPeripheralService(repo, peripheralAMQPHandler, client)

	// Set up the RabbitMQ connection string

	if err1 != nil {
		// log.Fatal("Error setting up RabbitMQ %v ", zap.error(err1))
		fmt.Fprintf(os.Stderr, "Error setting up RabbitMQ: %v\n", err1)
	}
	defer peripheralAMQPHandler.Close()

	// Initialize the HTTP handler with the Peripheral service and the RabbitMQ handler
	peripheralHTTPHandler := httphandler.NewPeripheralServiceHTTPHandler(router, peripheralService, peripheralAMQPHandler) // Adjusted for package

	// Set up the HTTP routes
	peripheralHTTPHandler.InitRouter()

	// Start the HTTP server
	err := router.Run(":" + cfg.Port)
	if err != nil {
		log.Fatal("Failed to run the server: %v", zap.Error(err))
	}

	docs.SwaggerInfo.Host = "localhost:" + cfg.Port
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
