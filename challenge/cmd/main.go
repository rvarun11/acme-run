package main

import (
	"github.com/CAS735-F23/macrun-teamvsl/challenge/config"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/adapters/handler/http"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/adapters/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/services"
	logger "github.com/CAS735-F23/macrun-teamvsl/challenge/log"
	"github.com/gin-gonic/gin"
)

var cfg *config.AppConfiguration = config.Config

func main() {
	logger.Info("Challenge Manager is starting...")

	// Initialize router
	router := gin.Default()
	router.Use(gin.Recovery())

	// Initialize postgres repository
	store := postgres.NewRepository(cfg.Postgres)

	// Initialize challenge service
	challengeSvc := services.NewChallengeService(store)
	challengeHandler := http.NewChallengeHandler(router, challengeSvc)
	challengeHandler.InitRouter()

	// Initialize badge service
	// badgeSvc := services.NewBadgeService(store)
	// badgeHandler := http.NewBadgeHandler(router, *badgeSvc)
	// badgeHandler.InitRouter()

	// Start Server
	router.Run(":" + cfg.Port)
}
