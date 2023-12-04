package main

import (
	"github.com/CAS735-F23/macrun-teamvsl/workout/config"
	amqpPrimary "github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/primary/amqp"
	http "github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/primary/http"
	amqpSecondary "github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/secondary/amqp"
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
	logger.Info("workout manager is starting...")

	// Initialize router
	router := gin.New()
	router.Use(gin.Recovery())

	// Initialize postgres repository
	store := postgres.NewRepository(cfg.Postgres)

	// Initialize clients
	peripheralClient := clients.NewPeripheralClient(cfg.PeripheralClient)
	userClient := clients.NewUserServiceClient(cfg.UserClient)

	// Initialize workout stats publisher
	workoutStatsWorkoutStatsPublisher := amqpSecondary.NewWorkoutStatsPublisher(cfg.RabbitMQ)

	// Initialize workout service
	workoutSvc := services.NewWorkoutService(store, peripheralClient, userClient, workoutStatsWorkoutStatsPublisher)
	workoutHandler := http.NewWorkoutHanlder(router, workoutSvc)
	workoutHandler.InitRouter()

	// Initialize shelter distance consumer
	shelterDistanceConsumer := amqpPrimary.NewShelterDistanceConsumer(cfg.RabbitMQ, workoutSvc)
	shelterDistanceConsumer.InitAMQP()

	// Initialize location consumer
	locationConsumer := amqpPrimary.NewLocationConsumer(cfg.RabbitMQ, workoutSvc)
	locationConsumer.InitAMQP()

	// Swagger support
	docs.SwaggerInfo.Host = "localhost:" + cfg.Port
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	router.Run(":" + cfg.Port)
}
