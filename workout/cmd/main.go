package main

import (
	"flag"
	"fmt"

	"github.com/CAS735-F23/macrun-teamvs_/workout/internal/adapters/handler"
	"github.com/CAS735-F23/macrun-teamvs_/workout/internal/adapters/repository"
	"github.com/CAS735-F23/macrun-teamvs_/workout/internal/core/services"
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

	InitRoutes()

}

func InitRoutes() {
	router := gin.Default()
	handler := handler.NewHTTPHandler(*svc)
	router.GET("/workouts", handler.ListWorkouts)
	router.POST("/workout", handler.StartWorkout)
	router.GET("/workouts/:id", handler.GetWorkout)
	// TODO: Implement when needed
	// router.POST("/player", handler.UpdatePlayer)
	router.Run(":8001")

}
