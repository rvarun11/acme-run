package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/CAS735-F23/macrun-teamvs_/hrm/internal/adapters/handler"
	"github.com/CAS735-F23/macrun-teamvs_/hrm/internal/adapters/repository"
	"github.com/CAS735-F23/macrun-teamvs_/hrm/internal/core/services"
	"github.com/gin-gonic/gin"
)

var (
	repo = flag.String("db", "postgres", "Database for storing messages")
	//    redisHost   = "localhost:6379"
	//    httpHandler *handler.HTTPHandler
	svc *services.HRMService
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
		svc = services.NewHRMService(store)
	}
	var wg sync.WaitGroup
	wg.Add(2)

	go InitRabbitMQ(&wg)
	go svc.SendHRM(&wg)
	InitRoutes()

	wg.Wait()
}

func InitRabbitMQ(wg *sync.WaitGroup) {

	defer wg.Done()

	handler.HRMWorkoutBinder(*svc)
}

func InitRoutes() {
	router := gin.Default()
	handler := handler.NewHTTPHandler(*svc)
	router.POST("/hrms", handler.ConnectHRM)
	// TODO: Implement when needed
	// router.PUT("/player", handler.UpdatePlayer)
	router.Run(":8004")

}
