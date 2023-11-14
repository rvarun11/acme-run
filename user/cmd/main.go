package main

import (
	"github.com/CAS735-F23/macrun-teamvsl/user/config"
	"github.com/CAS735-F23/macrun-teamvsl/user/internal/adapters/handler/http"
	"github.com/CAS735-F23/macrun-teamvsl/user/internal/adapters/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/user/internal/core/services"
	"github.com/CAS735-F23/macrun-teamvsl/user/logger"
	"github.com/gin-gonic/gin"
)

var cfg *config.AppConfiguration = config.Config

func main() {
	logger.Info("Challenge Manager is starting")

	// Initialize router
	router := gin.Default()
	router.Use(gin.Recovery())

	// Initialize postgres repository
	store := postgres.NewRepository(cfg.Postgres)

	// Initialize player service
	playerSvc := services.NewPlayerService(store)
	playerHandler := http.NewPlayerHandler(router, playerSvc)
	playerHandler.InitRouter()

	// Start Server
	router.Run(":" + cfg.Port)
}
