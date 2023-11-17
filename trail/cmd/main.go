package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/CAS735-F23/macrun-teamvsl/trail/config"
	httphandler "github.com/CAS735-F23/macrun-teamvsl/trail/internal/adapters/handler/http"
	"github.com/CAS735-F23/macrun-teamvsl/trail/internal/adapters/handler/peripheralhandler"
	repository "github.com/CAS735-F23/macrun-teamvsl/trail/internal/adapters/repository/memory"
	"github.com/CAS735-F23/macrun-teamvsl/trail/internal/adapters/repository/postgres"
	"github.com/CAS735-F23/macrun-teamvsl/trail/internal/core/services"
	log "github.com/CAS735-F23/macrun-teamvsl/trail/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var cfg *config.AppConfiguration = config.Config

func printTableData(db *sql.DB, tableName string) {
	// Check if the table is empty
	var count int
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&count)
	if err != nil {
		log.Fatal("Error counting rows in print", zap.Error(err))
	}

	// If the table is not empty, print its contents
	if count > 0 {
		rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
		if err != nil {
			// log.Fatalf("Error querying data from %s: %v", tableName, err)
		}
		defer rows.Close()

		cols, err := rows.Columns()
		if err != nil {
			// log.Fatalf("Error getting columns from %s: %v", tableName, err)
		}

		values := make([]interface{}, len(cols))
		pointers := make([]interface{}, len(cols))
		for i := range values {
			pointers[i] = &values[i]
		}

		for rows.Next() {
			err := rows.Scan(pointers...)
			if err != nil {
				// log.Fatal(err)
			}

			for i, colName := range cols {
				val := pointers[i].(*interface{})
				fmt.Printf("%s: %v ", colName, *val)
			}
			fmt.Println()
		}
		if err = rows.Err(); err != nil {
			// log.Fatal(err)
		}
	} else {
		fmt.Printf("The table %s is empty.\n", tableName)
	}
}

func initializeDB() {
	num, err := strconv.Atoi(cfg.Postgres.Port)
	if err != nil {
		// Handle the error if the string cannot be converted to an integer
		fmt.Println("Error:", err)
		return
	}
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, num, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DB_Name)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		// log.Fatal(err)
		fmt.Println("Error:", err)
	}
	defer db.Close()

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("database cant be connected", zap.Error(err))
	}
	fmt.Println("Successfully connected!")

	// Create a new schema
	_, err = db.Exec("CREATE SCHEMA IF NOT EXISTS traildb AUTHORIZATION guest")
	if err != nil {
		log.Fatal("jkjk", zap.Error(err))
	}
	fmt.Println("Schema created.")

	// Grant all privileges on the schema to the user
	_, err = db.Exec("GRANT ALL PRIVILEGES ON SCHEMA traildb TO guest")
	if err != nil {
		log.Fatal("database access cant be granted", zap.Error(err))
	}
	fmt.Println("Privileges granted to user 'guest'.")

	// Create tables
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS traildb.zone (
            zone_id UUID PRIMARY KEY,
            zone_name TEXT NOT NULL
        )`)
	if err != nil {
		log.Fatal("Error creating zone table", zap.Error(err))
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS traildb.trail (
            trail_id UUID PRIMARY KEY,
            trail_name TEXT NOT NULL,
            zone_id UUID REFERENCES traildb.zone(zone_id),
            start_longitude FLOAT,
            start_latitude FLOAT,
            end_longitude FLOAT,
            end_latitude FLOAT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`)
	if err != nil {
		log.Fatal("cannot inser to the trail table", zap.Error(err))
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS traildb.shelter (
            shelter_id UUID PRIMARY KEY,
            shelter_name TEXT NOT NULL,
            trail_id UUID REFERENCES traildb.trail(trail_id),
            shelter_availability BOOLEAN NOT NULL,
            longitude FLOAT,
            latitude FLOAT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`)
	if err != nil {
		log.Fatal("Error creating shelter table", zap.Error(err))
	}

	fmt.Println("Table created.")

	// Insert two shelters with distinct names and locations

	zones := []struct {
		ID   uuid.UUID
		Name string
	}{
		{uuid.New(), "McMaster Zone"},
		{uuid.New(), "Fortinos Zone"},
		{uuid.New(), "Westdale Zone"},
	}

	// Insert four trails with distinct names and start/end longitudes
	trails := []struct {
		ID       uuid.UUID
		ZoneID   uuid.UUID
		Name     string
		StartLon float64
		StartLat float64
		EndLon   float64
		EndLat   float64
	}{
		{uuid.New(), zones[0].ID, "Cedar Pass Trail", 40.1, 45.0, 42.3, 45.0},
		{uuid.New(), zones[0].ID, "Blue Ridge Path", 42.5, 45.0, 44.7, 45.0},
		{uuid.New(), zones[1].ID, "Redwood Walk", 45.2, 45.0, 47.8, 45.0},
		{uuid.New(), zones[2].ID, "Willow Way", 46.3, 45.0, 49.9, 45.0},
		{uuid.New(), zones[2].ID, "Willow2 Way", 46.3, 45.0, 49.9, 45.0},
		{uuid.New(), zones[2].ID, "Willow3 Way", 46.3, 45.0, 49.9, 45.0},
	}

	shelters := []struct {
		ID      uuid.UUID
		Name    string
		TrailID uuid.UUID
		Avail   bool
		Lon     float64
		Lat     float64
	}{
		{uuid.New(), "Westdale Shelter1", trails[0].ID, true, 45.0, 45.0},
		{uuid.New(), "Westdale Shelter2", trails[0].ID, true, 46.0, 46.0},
		{uuid.New(), "Fortinos Shelter", trails[2].ID, true, 46.0, 46.0},
		{uuid.New(), "McMaster Shelter", trails[3].ID, false, 46.0, 46.0},
	}

	// Insert initial data into zone table
	// Note: Adjust the UUIDs to your preference or generate them programmatically
	var zoneCount int
	err = db.QueryRow("SELECT COUNT(*) FROM traildb.zone").Scan(&zoneCount)
	if err != nil {
		log.Fatal("error counting zones: %v", zap.Error(err))
	}
	if zoneCount < 2 {
		for _, zone := range zones {
			_, err := db.Exec(`
            INSERT INTO traildb.zone (zone_id, zone_name) VALUES ($1, $2)
            ON CONFLICT (zone_id) DO NOTHING`,
				zone.ID, zone.Name)
			if err != nil {
				log.Fatal("Error inserting zone", zap.Error(err))
			}
		}
	}

	fmt.Println("inserted zones.")

	var trailCount int
	err = db.QueryRow("SELECT COUNT(*) FROM traildb.trail").Scan(&trailCount)
	if err != nil {
		log.Fatal("error counting trails: %v", zap.Error(err))
	}

	// Create a table within the schema
	if trailCount < len(trails) {
		for _, trail := range trails {
			_, err := db.Exec(`
            INSERT INTO traildb.trail (trail_id, trail_name, zone_id, start_longitude, start_latitude, end_longitude, end_latitude)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
            ON CONFLICT (trail_id) DO NOTHING`,
				trail.ID, trail.Name, trail.ZoneID, trail.StartLon, trail.StartLat, trail.EndLon, trail.EndLat)
			if err != nil {
				log.Fatal("Error inserting trail", zap.Error(err))
			}
		}
	}

	fmt.Println("inserted trails.")
	var shelterCount int
	err = db.QueryRow("SELECT COUNT(*) FROM traildb.shelter").Scan(&shelterCount)
	if err != nil {
		log.Fatal("error counting shelters: %v", zap.Error(err))
	}
	if shelterCount < len(shelters) {
		for _, shelter := range shelters {
			_, err := db.Exec(`
            INSERT INTO traildb.shelter (shelter_id, shelter_name, trail_id, shelter_availability, longitude, latitude)
            VALUES ($1, $2, $3, $4, $5, $6)
            ON CONFLICT (shelter_id) DO NOTHING`,
				shelter.ID, shelter.Name, shelter.TrailID, shelter.Avail, shelter.Lon, shelter.Lat)
			if err != nil {
				log.Fatal("Error inserting shelter", zap.Error(err))
			}
		}
	}

	// printTableData(db, "traildb.shelter")
	// printTableData(db, "traildb.trail")
}

func main() {
	log.Info("Trail Service is starting")

	// Initialize router
	router := gin.New()
	router.Use(gin.Recovery())

	// initialize a database for test purposes
	initializeDB()

	// Initialize the repository
	repo := repository.NewMemoryRepository()
	repoS := postgres.NewShelterRepository(cfg.Postgres)
	repoT := postgres.NewTrailRepository(cfg.Postgres)
	repoZ := postgres.NewZoneRepository(cfg.Postgres)

	// Initialize the trailManager service
	trailManagerService, _ := services.NewTrailManagerService(repo, repoT, repoS, repoZ)

	// Initialize the HTTP handler with the trail manager service and the RabbitMQ handler
	trailManagerHTTPHandler := httphandler.NewTrailManagerServiceHTTPHandler(router, trailManagerService) // Adjusted for package

	// Set up the HTTP routes
	trailManagerHTTPHandler.InitRouter()
	phandler := peripheralhandler.NewAMQPHandler(trailManagerService)
	phandler.InitAMQP()

	// Start the HTTP server
	err := router.Run(":" + cfg.Port)
	if err != nil {
		log.Fatal("Failed to run the server", zap.Error(err))
	}
}
