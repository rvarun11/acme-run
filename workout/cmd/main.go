package main

import (
	"github.com/CAS735-F23/macrun-teamvsl/workout/config"
	amqphandler "github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/primary/handler/amqp"
	http "github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/primary/handler/http"
	amqpsecondaryadapter "github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/secondary/amqp"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/secondary/clients"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/secondary/repository/postgres"
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

	workoutAMQPSecondaryHandler, _ := amqpsecondaryadapter.NewAMQPPublisher()

	// Initialize Workout service
	workoutSvc := services.NewWorkoutService(store, peripheralDeviceClient, user, workoutAMQPSecondaryHandler)

	workoutHTTPHandler := http.NewWorkoutHTTPHandler(router, workoutSvc)
	workoutHTTPHandler.InitRouter()

	workoutAMQPHandler := amqphandler.NewAMQPHandler(workoutSvc)

	go workoutAMQPHandler.InitAMQP()
	// Start Server
	router.Run(":" + cfg.Port)
}
