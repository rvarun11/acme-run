package main

import (
	"github.com/CAS735-F23/macrun-teamvsl/user/config"
	"github.com/CAS735-F23/macrun-teamvsl/user/docs"
	"github.com/CAS735-F23/macrun-teamvsl/user/internal/adapters/primary/http"
	"github.com/CAS735-F23/macrun-teamvsl/user/internal/adapters/secondary/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/user/internal/core/services"
	logger "github.com/CAS735-F23/macrun-teamvsl/user/log"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var cfg *config.AppConfiguration = config.Config

// @title User Manager API
// @version 1.0
// @description This provides a description of API endpoints for the Player Manager

// @contact.name Varun Rajput
// @contact.url    https://github.com/rvarun11
// @contact.email rajpuv2@mcmaster.ca

// @query.collection.format multi
func main() {
	logger.Info("user manager is starting...")

	// Initialize router
	router := gin.Default()
	router.Use(gin.Recovery())

	// Initialize postgres repository
	store := postgres.NewRepository(cfg.Postgres)

	// Initialize player service
	playerSvc := services.NewPlayerService(store)
	playerHandler := http.NewPlayerHandler(router, playerSvc)
	playerHandler.InitRouter()

	// Swagger Support
	docs.SwaggerInfo.Host = "localhost:" + cfg.Port
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start Server
	router.Run(":" + cfg.Port)
}
