package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/trail/config"
	httphandler "github.com/CAS735-F23/macrun-teamvsl/trail/internal/adapters/handler/http"
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

func createTables(db *sql.DB, dbName string) error {
	createTablesSQL := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s.trail (
		trail_id UUID PRIMARY KEY,
		trail_name TEXT NOT NULL,
		zone_id INT NOT NULL,
		start_longitude FLOAT,
		start_latitude FLOAT,
		end_longitude FLOAT,
		end_latitude FLOAT,
		shelter_id UUID,
		created_at TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS %s.shelter (
		shelter_id UUID PRIMARY KEY,
		shelter_name TEXT NOT NULL,
		zone_id INT NOT NULL,
		shelter_availability BOOLEAN,
		longitude FLOAT,
		latitude FLOAT
	);
	`, dbName, dbName)

	_, err := db.Exec(createTablesSQL)
	if err != nil {
		return err
	}

	return nil
}

func createAndInsertData(db *sql.DB) error {
	// Define SQL statements for table creation
	createTrailTableSQL := `
	CREATE TABLE IF NOT EXISTS trail (
		trail_id UUID PRIMARY KEY,
		trail_name TEXT,
		zone_id INT,
		start_longitude DOUBLE PRECISION,
		start_latitude DOUBLE PRECISION,
		end_longitude DOUBLE PRECISION,
		end_latitude DOUBLE PRECISION,
		shelter_id UUID,
		created_at TIMESTAMP
	);
	`

	createShelterTableSQL := `
	CREATE TABLE IF NOT EXISTS shelter (
		shelter_id UUID PRIMARY KEY,
		shelter_name TEXT,
		zone_id INT,
		shelter_availability BOOLEAN,
		longitude DOUBLE PRECISION,
		latitude DOUBLE PRECISION
	);
	`

	// Execute table creation SQL statements
	if _, err := db.Exec(createTrailTableSQL); err != nil {
		return err
	}

	if _, err := db.Exec(createShelterTableSQL); err != nil {
		return err
	}

	// Check if Trail table is empty and insert sample data
	var trailCount int
	err := db.QueryRow("SELECT COUNT(*) FROM trail").Scan(&trailCount)
	if err != nil {
		return fmt.Errorf("failed to count trails: %w", err)
	}

	if trailCount == 0 {
		for i := 1; i <= 3; i++ {
			trailID := uuid.New()
			trailName := fmt.Sprintf("Trail %d", i)
			zoneID := i
			startLongitude := float64(i)
			startLatitude := float64(i)
			endLongitude := float64(i + 1)
			endLatitude := float64(i + 1)
			shelterID := uuid.New()
			createdAt := time.Now()

			_, err := db.Exec(`
				INSERT INTO trail (trail_id, trail_name, zone_id, start_longitude, start_latitude, end_longitude, end_latitude, shelter_id, created_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`, trailID, trailName, zoneID, startLongitude, startLatitude, endLongitude, endLatitude, shelterID, createdAt)

			if err != nil {
				return err
			}
		}
	}

	// Check if the shelter table is empty and insert sample data if it is
	var shelterCount int
	err = db.QueryRow("SELECT COUNT(*) FROM shelter").Scan(&shelterCount)
	if err != nil {
		log.Error("cant count shelter")
		return err // Handle the error appropriately
	}

	if shelterCount == 0 {
		for i := 1; i <= 3; i++ {
			shelterID := uuid.New()
			shelterName := fmt.Sprintf("Shelter %d", i)
			zoneID := i
			shelterAvailability := true
			longitude := float64(i)
			latitude := float64(i)

			_, err := db.Exec(`
				INSERT INTO shelter (shelter_id, shelter_name, zone_id, shelter_availability, longitude, latitude)
				VALUES ($1, $2, $3, $4, $5, $6)
			`, shelterID, shelterName, zoneID, shelterAvailability, longitude, latitude)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

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

	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS traildb.trail (
				trail_id UUID PRIMARY KEY,
				trail_name TEXT NOT NULL,
				zone_id INT NOT NULL,
				start_longitude FLOAT,
				start_latitude FLOAT,
				end_longitude FLOAT,
				end_latitude FLOAT,
				shelter_id UUID,
				created_at TIMESTAMP
			)`)
	if err != nil {
		log.Fatal("cannot inser to the trail table", zap.Error(err))
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS traildb.shelter (
        shelter_id UUID PRIMARY KEY,
        shelter_name TEXT NOT NULL,
        zone_id INT NOT NULL,
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
	shelters := []struct {
		ID   uuid.UUID
		Name string
		Lon  float64
		Lat  float64
	}{
		{uuid.New(), "Maplewood Shelter", 45.0, 45.0},
		{uuid.New(), "Riverside Shelter", 46.0, 46.0},
	}

	var shelterCount int
	err = db.QueryRow("SELECT COUNT(*) FROM traildb.shelter").Scan(&shelterCount)
	if err != nil {
		log.Fatal("error counting shelters: %v", zap.Error(err))
	}
	if shelterCount < 2 {

		for _, shelter := range shelters {
			_, err := db.Exec(`
            INSERT INTO traildb.shelter (shelter_id, shelter_name, zone_id, shelter_availability, longitude, latitude)
            VALUES ($1, $2, $3, $4, $5, $6)`,
				shelter.ID, shelter.Name, 1, true, shelter.Lon, shelter.Lat)
			if err != nil {
				log.Fatal("error inserting shelter", zap.Error(err))
			}
		}
	}

	// Insert four trails with distinct names and start/end longitudes
	trails := []struct {
		Name      string
		StartLon  float64
		StartLat  float64
		EndLon    float64
		EndLat    float64
		ShelterID uuid.UUID
	}{
		{"Cedar Pass Trail", 40.1, 45.0, 42.3, 45.0, uuid.Nil},
		{"Blue Ridge Path", 42.5, 45.0, 44.7, 45.0, uuid.Nil},
		{"Redwood Walk", 45.2, 45.0, 47.8, 45.0, shelters[0].ID},
		{"Willow Way", 46.3, 45.0, 49.9, 45.0, shelters[1].ID},
	}

	var trailCount int
	err = db.QueryRow("SELECT COUNT(*) FROM traildb.trail").Scan(&trailCount)
	if err != nil {
		log.Fatal("error counting trails: %v", zap.Error(err))
	}

	initCounter := 0
	// Create a table within the schema
	if trailCount < 4 {
		for _, trail := range trails {

			_, err := db.Exec(`
            INSERT INTO traildb.trail (trail_id, trail_name, zone_id, start_longitude, start_latitude, end_longitude, end_latitude, shelter_id, created_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
				uuid.New(), trail.Name, initCounter/2+1, trail.StartLon, trail.StartLat, trail.EndLon, trail.EndLat, trail.ShelterID, time.Now())
			if err != nil {
				log.Fatal("error inserting trail: %v", zap.Error(err))
			}
			initCounter += 1
		}
	}

	printTableData(db, "traildb.shelter")
	printTableData(db, "traildb.trail")
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

	// Initialize the trailManager service
	trailManagerService, _ := services.NewTrailManagerService(repo, repoT, repoS)
	// trailManagerService.

	// Connect to the default database to perform administrative tasks

	// Set up the RabbitMQ connection string
	// amqpURL := "amqp://" + cfg.RabbitMQ.User + ":" +
	// 	cfg.RabbitMQ.Password + "@" + cfg.RabbitMQ.Host + ":" + cfg.RabbitMQ.Port + "/"

	// Initialize the RabbitMQ handler with the Peripheral service and the AMQP URL
	// peripheralAMQPHandler, err1 := handler.NewRabbitMQHandler(peripheralService, amqpURL) // Adjusted for package
	// if err1 != nil {
	// 	// log.Fatal("Error setting up RabbitMQ %v ", zap.error(err1))
	// 	fmt.Fprintf(os.Stderr, "Error setting up RabbitMQ: %v\n", err1)
	// }
	// defer peripheralAMQPHandler.Close()

	// Initialize the HTTP handler with the trail manager service and the RabbitMQ handler
	trailManagerHTTPHandler := httphandler.NewTrailManagerServiceHTTPHandler(router, trailManagerService) // Adjusted for package

	// Set up the HTTP routes
	trailManagerHTTPHandler.InitRouter()

	// Start the HTTP server
	err := router.Run(":" + cfg.Port)
	if err != nil {
		// log.Fatal("Failed to run the server: %v", err)
	}
}
