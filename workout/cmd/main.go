package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/handler"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/services"
	"github.com/gin-gonic/gin"
)

var svc *services.WorkoutService

func main() {
	flag.Parse()

	fmt.Println("Workout service is running")
	store := postgres.NewRepository()
	svc = services.NewWorkoutService(store)

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
	router.POST("/workout", handler.StartWorkout)
	router.PUT("/workout", handler.StopWorkout)
	router.GET("/workoutOptions", handler.GetWorkoutOptions)
	router.POST("/workoutOptions", handler.StartWorkoutOption)
	router.PUT("/workoutOptions", handler.StopWorkoutOption)

	router.GET("workout/distance", handler.GetDistance)
	router.GET("workout/shelters", handler.GetShelters)
	router.GET("workout/escapes", handler.GetEscapes)
	router.GET("workout/fights", handler.GetFights)

	router.Run(":8001")
}

// TODO: Handle service configurations properly
type Config struct {
	RABBITMQ_URL string
}

func NewConfig() Config {
	return Config{
		RABBITMQ_URL: GetEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
	}
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
