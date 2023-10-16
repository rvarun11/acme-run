package main

import (
	"flag"
	"fmt"

	"github.com/CAS735-F23/macrun-teamvs_/hrm/adapters/repository"
	"github.com/CAS735-F23/macrun-teamvs_/hrm/internal/adapters/handler"
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

	InitRoutes()

}

func InitRoutes() {
	router := gin.Default()
	handler := handler.NewHTTPHandler(*svc)
	router.GET("/hrmss", handler.ListHRM)
	router.POST("/hrm", handler.CreateHRM)
	router.GET("/hrmss/:id", handler.GetHRM)
	// TODO: Implement when needed
	// router.PUT("/player", handler.UpdatePlayer)
	router.Run(":8000")

}
