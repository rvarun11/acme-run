package main

import (
	"fmt"

	"github.com/CAS735-F23/macrun-teamvsl/player/internal/adapters/handler"
	"github.com/CAS735-F23/macrun-teamvsl/player/internal/adapters/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/player/internal/core/services"
	"github.com/gin-gonic/gin"
)

var svc *services.PlayerService

func main() {
	fmt.Println("Player Service is running")
	store := postgres.NewRepository()
	svc = services.NewPlayerService(store)
	InitRoutes()
}

func InitRoutes() {
	router := gin.Default()
	handler := handler.NewHTTPHandler(*svc)

	// Player Routes
	router.POST("/players", handler.RegisterPlayer)
	router.GET("/players/:id", handler.GetPlayer)
	router.PUT("/players", handler.UpdatePlayer)
	router.GET("/players", handler.ListPlayers)

	// Admin Routes
	// router.POST("/player", handler.CreateAdmin)
	// router.GET("/players/:id", handler.GetAdmin)
	// router.PUT("/player", handler.UpdateAdmin)
	// router.GET("/players", handler.ListPlayers)
	router.Run(":8000")
}
