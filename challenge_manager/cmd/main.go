package main

import (
	"github.com/CAS735-F23/macrun-teamvsl/challenge_manager/internal/adapters/handler/http"
	"github.com/CAS735-F23/macrun-teamvsl/challenge_manager/internal/adapters/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/challenge_manager/internal/core/services"
	log "github.com/CAS735-F23/macrun-teamvsl/challenge_manager/log"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Info("Challenge Manager is starting")
	// Initialize router
	router := gin.New()
	router.Use(gin.Recovery())

	// Initialize postgres repository
	store := postgres.NewRepository()

	// Initialize Challenge Manager
	challengeSvc := services.NewChallengeService(store)
	challengeHandler := http.NewChallengeHandler(router, *challengeSvc)
	challengeHandler.InitRouter()

	// Initialize Badge Manager
	// badgeSvc := services.NewBadgeService(store)
	// badgeHandler := http.NewBadgeHandler(router, *badgeSvc)
	// badgeHandler.InitRouter()

	// Start Server
	router.Run(":8001")
}
