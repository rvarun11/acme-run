package main

import (
	"github.com/CAS735-F23/macrun-teamvsl/workout/config"
	amqphandler "github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/primary/amqp"
	http "github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/primary/http"
	amqpsecondaryadapter "github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/secondary/amqp"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/secondary/clients"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/secondary/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/services"
	logger "github.com/CAS735-F23/macrun-teamvsl/workout/log"
	"github.com/gin-gonic/gin"

	"github.com/CAS735-F23/macrun-teamvsl/workout/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var cfg *config.AppConfiguration = config.Config

func main() {
	logger.Info("Workout Service is starting...")

	// Initialize router
	router := gin.New()
	router.Use(gin.Recovery())

	// Initialize postgres repository
	store := postgres.NewRepository(cfg.Postgres)

	// Initialize Clients
	peripheralDeviceClient := clients.NewPeripheralDeviceClient()
	user := clients.NewUserServiceClient()

	workoutAMQPSecondaryHandler := amqpsecondaryadapter.NewPublisher(cfg.RabbitMQ)

	// Initialize Workout service
	workoutSvc := services.NewWorkoutService(store, peripheralDeviceClient, user, workoutAMQPSecondaryHandler)

	workoutHTTPHandler := http.NewWorkoutHTTPHandler(router, workoutSvc)
	workoutHTTPHandler.InitRouter()

	workoutAMQPHandler := amqphandler.NewAMQPHandler(workoutSvc)

	go workoutAMQPHandler.InitAMQP()
	// Start Server

	docs.SwaggerInfo.Host = "localhost:" + cfg.Port
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":" + cfg.Port)
}
