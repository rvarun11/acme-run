package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/trail/config"
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
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM trail.%s", tableName)).Scan(&count)
	if err != nil {
		// log.Fatalf("Error counting rows in %s: %v", tableName, err)
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
		log.Fatal("jkjk", zap.Error(err))
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
		log.Fatal("jkjk", zap.Error(err))
	}
	fmt.Println("Privileges granted to user 'guest'.")

	// Create a table within the schema
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
		log.Fatal("jkjk", zap.Error(err))
	}
	fmt.Println("Table created.")
	_, err = db.Exec(`
        INSERT INTO traildb.trail (trail_id, trail_name, zone_id, start_longitude, start_latitude, end_longitude, end_latitude, shelter_id, created_at)
        VALUES
            ('123e4567-e89b-12d3-a456-426614174000', 'Appalachian Trail', 1, -81.564, 36.132, -81.811, 36.273, '123e4567-e89b-12d3-a456-426614174001', $1),
            ('123e4567-e89b-12d3-a456-426614174002', 'Pacific Crest Trail', 2, -121.236, 36.578, -121.489, 36.845, '123e4567-e89b-12d3-a456-426614174003', $1)
        ON CONFLICT (trail_id) DO NOTHING`, time.Now())
	if err != nil {
		log.Fatal("jkjk", zap.Error(err))
	}

	printTableData(db, "shelter")
	printTableData(db, "trail")
}

func main() {
	log.Info("Trail Service is starting")

	// Initialize router
	router := gin.New()
	router.Use(gin.Recovery())

	// initialize a database for test purposes
	initializeDB()

	// Initialize the repository
	// repo := repository.NewMemoryRepository()

	// Initialize the trailManager service
	// trailManagerService := services.NewTrailManagerService(repo)
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

	// // Initialize the HTTP handler with the Peripheral service and the RabbitMQ handler
	// peripheralHTTPHandler := handler.NewPeripheralServiceHTTPHandler(router, peripheralService, peripheralAMQPHandler) // Adjusted for package

	// // Set up the HTTP routes
	// peripheralHTTPHandler.InitRouter()

	// // Start the HTTP server
	// err := router.Run(":" + cfg.Port)
	// if err != nil {
	// 	// log.Fatal("Failed to run the server: %v", err)
	// }
}
