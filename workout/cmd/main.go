package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/handler"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/repository"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/services"
	"github.com/gin-gonic/gin"
)

var (
	repo = flag.String("db", "postgres", "Database for storing messages")
	//    redisHost   = "localhost:6379"
	//    httpHandler *handler.HTTPHandler
	svc *services.WorkoutService
)

func main() {
	flag.Parse()

	fmt.Printf("Application running using %s\n", *repo)
	switch *repo {
	// note: we can have other repositories like redis, mysql, etc
	//    case "redis":
	//        store := repository.NewMessengerRedisRepository(redisHost)
	//        svc = services.NewMessengerService(store)
	default:
		store := repository.NewMemoryRepository()
		svc = services.NewWorkoutService(store)
	}
	var wg sync.WaitGroup
	wg.Add(1)

	go InitRabbitMQ(&wg)
	InitRoutes()

	wg.Wait()
}

func InitRabbitMQ(wg *sync.WaitGroup) {
	defer wg.Done()
	cfg := NewConfig()
	handler.HRMSubscriber(*svc, cfg.RABBITMQ_URL)
}

func InitRoutes() {

	router := gin.Default()
	handler := handler.NewHTTPHandler(*svc)
	router.GET("/workouts", handler.ListWorkouts)
	router.POST("/workout", handler.StartWorkout)
	router.PUT("/workout", handler.StopWorkout)
	router.GET("/workouts/:id", handler.GetWorkout)
	// TODO: Implement when needed
	// router.POST("/player", handler.UpdatePlayer)
	router.Run(":8001")

}

// TODO: Handle service configurations properly
type Config struct {
	RABBITMQ_URL string
}

func NewConfig() Config {
	return Config{
		RABBITMQ_URL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
