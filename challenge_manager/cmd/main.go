package main

import (
	"fmt"

	"github.com/CAS735-F23/macrun-teamvsl/challenge_manager/internal/adapters/handler/http"
	"github.com/CAS735-F23/macrun-teamvsl/challenge_manager/internal/adapters/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/challenge_manager/internal/core/services"
	"github.com/gin-gonic/gin"
)

var svc *services.ChallengeService

func main() {
	fmt.Println("Challenge Service is running")
	store := postgres.NewRepository()
	svc = services.NewChallengeService(store)
	InitRoutes()
}

func InitRoutes() {
	router := gin.Default()
	handler := http.NewHandler(*svc)

	// Challenge Routes
	router.POST("/challenges", handler.CreateChallenge)
	router.GET("/challenges/:id", handler.GetChallengeByID)
	router.PUT("/challenges", handler.UpdateChallenge)
	router.GET("/challenges", handler.ListChallenges)
	// router.GET("/challenges", handler.ListActiveChallenges)

	// Badge Routes
	// router.POST("/player", handler.CreateAdmin)
	// router.GET("/players/:id", handler.GetAdmin)
	// router.PUT("/player", handler.UpdateAdmin)
	// router.GET("/players", handler.ListPlayers)
	router.Run(":8001")
}
