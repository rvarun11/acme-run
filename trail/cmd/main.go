package main

import (
	"fmt"
	"os"

	"github.com/CAS735-F23/macrun-teamvsl/trail/config"
	"github.com/CAS735-F23/macrun-teamvsl/trail/internal/adapters/handler"
	"github.com/CAS735-F23/macrun-teamvsl/trail/internal/adapters/repository"
	"github.com/CAS735-F23/macrun-teamvsl/trail/internal/core/services"
	log "github.com/CAS735-F23/macrun-teamvsl/trail/log"
	"github.com/gin-gonic/gin"
)

var cfg *config.AppConfiguration = config.Config

func main() {
	log.Info("Peripheral Service is starting")

	// Initialize router
	router := gin.New()
	router.Use(gin.Recovery())

	// Initialize the repository
	repo := repository.NewMemoryRepository()

	// Initialize the Peripheral service
	trailService := services.NewTrailManagerService(repo)

	// Set up the RabbitMQ connection string
	amqpURL := "amqp://" + cfg.RabbitMQ.User + ":" +
		cfg.RabbitMQ.Password + "@" + cfg.RabbitMQ.Host + ":" + cfg.RabbitMQ.Port + "/"

	// Initialize the RabbitMQ handler with the Peripheral service and the AMQP URL
	trailAMQPHandler, err1 := handler.NewRabbitMQHandler(trailService, amqpURL) // Adjusted for package
	if err1 != nil {
		// log.Fatal("Error setting up RabbitMQ %v ", zap.error(err1))
		fmt.Fprintf(os.Stderr, "Error setting up RabbitMQ: %v\n", err1)
	}
	defer trailAMQPHandler.Close()

	// Initialize the HTTP handler with the Peripheral service and the RabbitMQ handler
	trailHTTPHandler := handler.NewTrailServiceHTTPHandler(router, trailService, trailAMQPHandler) // Adjusted for package

	// Set up the HTTP routes
	trailHTTPHandler.InitRouter()

	// Start the HTTP server
	err := router.Run(":" + cfg.Port)
	if err != nil {
		// log.Fatal("Failed to run the server: %v", err)
	}
}
