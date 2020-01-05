package postgresdriver

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	sq "github.com/Masterminds/squirrel"

	"github.com/kolya59/easy_normalization/pkg/car"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var db *sql.DB

// Init database
func InitDatabaseConnection(host string, port string, user string, password string, name string) (err error) {
	// Open connection
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("could not open database connection: %v", err)
	}
	// Test connection
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("could not connect to database: %v", err)
	}
	return
}

// Init database structure
func InitDatabaseStructure() (err error) {
	// Get data from script
	path := "./script.sql"
	scriptFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	script := string(scriptFile)

	// Execute script
	_, err = db.Exec(script)
	if err != nil {
		return err
	}
	return nil
}

// Close db connection
func CloseConnection() (err error) {
	return db.Close()
}

// Send data to DB
func SaveCars(cars []car.Car) error {
	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare queries
	queries := make(map[string]sq.InsertBuilder)
	queries["insertEngineQuery"] = sq.Insert("engines")
	queries["insertTransmissionQuery"] = sq.Insert("transmissions")
	queries["insertBrandQuery"] = sq.Insert("brands")
	queries["insertWheelQuery"] = sq.Insert("wheels")
	queries["insertCarQuery"] = sq.Insert("cars").Columns("model", "engine", "transmission", "brand", "wheel", "price")

	// Bind arguments to queries
	for _, c := range cars {
		queries["insertEngineQuery"].Values(c.EngineModel, c.EnginePower, c.EngineVolume, c.EngineType)
		queries["insertTransmissionQuery"].Values(c.TransmissionModel, c.TransmissionType, c.TransmissionGearsNumber)
		queries["insertBrandQuery"].Values(c.BrandName, c.BrandCreatorCountry)
		queries["insertWheelQuery"].Values(c.WheelModel, c.WheelRadius, c.WheelColor)
		queries["insertCarQuery"].Values(c.Model, c.EngineModel, c.TransmissionModel, c.BrandName, c.WheelModel, c.Price)
	}

	// Execute queries
	for _, query := range queries {
		convertedQuery, _, err := query.ToSql()
		if err != nil {
			return err
		}
		convertedQuery += " ON CONFLICT DO NOTHING"
		if _, err = tx.Exec(convertedQuery); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error().Err(rbErr).Msg("Failed to rollback transaction")
			}
			return err
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Error().Err(rbErr).Msg("Failed to rollback transaction")
		}
		return err
	}

	return nil
}
