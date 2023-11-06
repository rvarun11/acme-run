package main

import (
	"github.com/CAS735-F23/macrun-teamvsl/workout/config"
	http "github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/handler/http"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/adapters/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/services"
	log "github.com/CAS735-F23/macrun-teamvsl/workout/log"
	"github.com/gin-gonic/gin"
)

var cfg *config.AppConfiguration = config.Config

func main() {
	log.Info("Workout Service is starting")

	// Initialize router
	router := gin.New()
	router.Use(gin.Recovery())

	// Initialize postgres repository
	store := postgres.NewRepository(cfg.Postgres)

	// Initialize Workout service
	workoutSvc := services.NewWorkoutService(store)
	workoutHTTPHandler := http.NewWorkoutHTTPHandler(router, workoutSvc)
	workoutHTTPHandler.InitRouter()

	// Start Server
	router.Run(":" + cfg.Port)
}
