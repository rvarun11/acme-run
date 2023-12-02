package main

import (
	"database/sql"
	"fmt"

	"github.com/CAS735-F23/macrun-teamvsl/zone/config"
	httphandler "github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/primary/http"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/primary/peripheralhandler"
	repository "github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/secondary/repository/memory"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/secondary/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/adapters/secondary/workouthandler"
	"github.com/CAS735-F23/macrun-teamvsl/zone/internal/core/services"
	log "github.com/CAS735-F23/macrun-teamvsl/zone/log"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var cfg *config.AppConfiguration = config.Config

// ensureDatabaseExists checks for the existence of the database and creates it if it doesn't exist
func ensureDatabaseExists(cfg *config.Postgres) error {
	// Connection string without the database name
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password)

	// Open a connection to the database server
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	// Check if the database exists
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname=$1);`
	err = db.QueryRow(query, cfg.DB_Name).Scan(&exists)
	if err != nil {
		return err
	}

	// If the database does not exist, create it
	if !exists {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s;", cfg.DB_Name))
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	log.Info("Trail Service is starting")

	// Initialize router
	router := gin.New()
	router.Use(gin.Recovery())

	// initialize a database for test purposes
	// initializeDB()
	err := ensureDatabaseExists(cfg.Postgres)

	// Initialize the repository
	repo := repository.NewMemoryRepository()
	repoS := postgres.NewShelterRepository(cfg.Postgres)
	repoT := postgres.NewTrailRepository(cfg.Postgres)
	repoZ := postgres.NewZoneRepository(cfg.Postgres)

	// Initialize the publisher for the zone managaer to send the shelter info to the queue
	workoutDataHandler, _ := workouthandler.NewAMQPPublisher()

	// Initialize the ZoneManager service
	ZoneService, _ := services.NewZoneService(repo, repoT, repoS, repoZ, workoutDataHandler)

	// Initialize the HTTP handler with the trail manager service and the RabbitMQ handler
	ZoneManagerHTTPHandler := httphandler.NewZoneServiceHTTPHandler(router, ZoneService) // Adjusted for package

	// Set up the HTTP routes
	ZoneManagerHTTPHandler.InitRouter()
	phandler := peripheralhandler.NewAMQPHandler(ZoneService)
	phandler.InitAMQP()

	// Start the HTTP server
	err = router.Run(":" + cfg.Port)
	if err != nil {
		log.Fatal("Failed to run the server", zap.Error(err))
	}
}
