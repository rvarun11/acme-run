package main

import (
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/config"
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/adapters/handler"
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/adapters/repository/inmemory"
	"github.com/CAS735-F23/macrun-teamvsl/peripheral/internal/core/services"
	log "github.com/CAS735-F23/macrun-teamvsl/peripheral/log"
	"github.com/gin-gonic/gin"
)

var cfg *config.AppConfiguration = config.Config

func main() {
	log.Info("Peripheral Service is starting")

	// Initialize router
	router := gin.New()
	router.Use(gin.Recovery())

	// Initialize the repository
	repo := inmemory.NewRepository()

	// Initialize the Peripheral service
	peripheralService := services.NewPeripheralService(repo)

	// Set up the RabbitMQ connection string
	amqpURL := "amqp://" + cfg.RabbitMQ.User + ":" +
		cfg.RabbitMQ.Password + "@" + cfg.RabbitMQ.Host + ":" + cfg.RabbitMQ.Port + "/"

	// Initialize the RabbitMQ handler with the Peripheral service and the AMQP URL
	peripheralAMQPHandler := handler.NewRabbitMQHandler(peripheralService, amqpURL) // Adjusted for package

	// Initialize the HTTP handler with the Peripheral service and the RabbitMQ handler
	peripheralHTTPHandler := handler.NewHTTPHandler(router, peripheralService, peripheralAMQPHandler) // Adjusted for package

	// Set up the HTTP routes
	peripheralHTTPHandler.InitRouter()

	// Start the HTTP server
	err := router.Run(":" + cfg.Port)
	if err != nil {
		log.Fatalf("Failed to run the server: %v", err)
	}
}

// package main

// import (
// 	"flag"
// 	"fmt"
// 	"os"
// 	"sync"

// 	"github.com/CAS735-F23/macrun-teamvsl/hrm/internal/adapters/handler"
// 	"github.com/CAS735-F23/macrun-teamvsl/hrm/internal/adapters/repository"
// 	"github.com/CAS735-F23/macrun-teamvsl/hrm/internal/core/services"
// 	"github.com/gin-gonic/gin"
// )

// var (
// 	repo = flag.String("db", "postgres", "Database for storing messages")
// 	//    redisHost   = "localhost:6379"
// 	//    httpHandler *handler.HTTPHandler
// 	svc *services.HRMService
// )

// func main() {
// 	// flag.Parse()

// 	fmt.Printf("Application running using %s\n", *repo)
// 	switch *repo {
// 	// note: we can have other repositories like redis, mysql, etc
// 	//    case "redis":
// 	//        store := repository.NewMessengerRedisRepository(redisHost)
// 	//        svc = services.NewMessengerService(store)
// 	default:
// 		store := repository.NewMemoryRepository()
// 		svc = services.NewHRMService(store)
// 	}
// 	var wg sync.WaitGroup
// 	wg.Add(2)

// 	go InitRabbitMQ(&wg)
// 	go svc.SendHRM(&wg)
// 	InitRoutes()

// 	wg.Wait()
// }

// func InitRabbitMQ(wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	cfg := NewConfig()
// 	handler.HRMWorkoutBinder(*svc, cfg.RABBITMQ_URL)
// }

// func InitRoutes() {
// 	router := gin.Default()
// 	handler := handler.NewHTTPHandler(*svc)
// 	router.POST("/hrms", handler.ConnectHRM)
// 	// TODO: Implement when needed
// 	// router.PUT("/player", handler.UpdatePlayer)
// 	router.Run(":8004")

// }

// // TODO: Handle service configurations properly
// type Config struct {
// 	RABBITMQ_URL string
// }

// func NewConfig() Config {
// 	return Config{
// 		RABBITMQ_URL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
// 	}
// }

// func getEnv(key, defaultValue string) string {
// 	if value, exists := os.LookupEnv(key); exists {
// 		return value
// 	}
// 	return defaultValue
// }
