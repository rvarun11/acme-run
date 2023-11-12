package main

import (
	"github.com/CAS735-F23/macrun-teamvsl/workout/config"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/clients"
	amqphandler "github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/handler/amqp"
	http "github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/handler/http"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/services"
	log "github.com/CAS735-F23/macrun-teamvsl/workout/log"
	"github.com/gin-gonic/gin"
)

var cfg *config.AppConfiguration = config.Config

func main() {
	log.Info("Workout Service is starting")

	// Initialize router
	router := gin.New()
	router.Use(gin.Recovery())

	// Initialize postgres repository
	store := postgres.NewRepository(cfg.Postgres)
	peripheralDeviceClient := clients.NewPeripheralDeviceClient()
	user := clients.NewUserServiceClient()

	// Initialize Workout service
	workoutSvc := services.NewWorkoutService(store, peripheralDeviceClient, user)
	workoutHTTPHandler := http.NewWorkoutHTTPHandler(router, workoutSvc)
	workoutHTTPHandler.InitRouter()

	workoutAMQPHandler := amqphandler.NewAMQPHandler(workoutSvc)
	workoutAMQPHandler.InitAMQP()

	// Start Server
	router.Run(":" + cfg.Port)
}
