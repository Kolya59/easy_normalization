package postgresdriver

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	sq "github.com/Masterminds/squirrel"

	pb "github.com/kolya59/easy_normalization/proto"

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
func SaveCars(cars []pb.Car) error {
	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare queries
	queries := make(map[string]sq.InsertBuilder)
	queries["insertEngineQuery"] = sq.Insert("engines").
		Columns("engine_model", "engine_power", "engine_volume", "engine_type")
	queries["insertTransmissionQuery"] = sq.Insert("transmissions").
		Columns("transmission_model", "transmission_type", "transmission_gears_number")
	queries["insertBrandQuery"] = sq.Insert("brands").
		Columns("brand_name", "brand_creator_country")
	queries["insertWheelQuery"] = sq.Insert("wheels").
		Columns("wheel_radius", "wheel_color", "wheel_model")
	queries["insertCarQuery"] = sq.Insert("cars").
		Columns("model", "engine", "transmission", "brand", "wheel", "price")

	// Bind arguments to queries
	for _, c := range cars {
		queries["insertEngineQuery"] = queries["insertEngineQuery"].Values(c.EngineModel, c.EnginePower, c.EngineVolume, c.EngineType)
		queries["insertTransmissionQuery"] = queries["insertTransmissionQuery"].Values(c.TransmissionModel, c.TransmissionType, c.TransmissionGearsNumber)
		queries["insertBrandQuery"] = queries["insertBrandQuery"].Values(c.BrandName, c.BrandCreatorCountry)
		queries["insertWheelQuery"] = queries["insertWheelQuery"].Values(c.WheelModel, c.WheelRadius, c.WheelColor)
		queries["insertCarQuery"] = queries["insertCarQuery"].Values(c.Model, c.EngineModel, c.TransmissionModel, c.BrandName, c.WheelModel, c.Price)
	}

	// Execute queries
	for _, query := range queries {
		query = query.Suffix("ON CONFLICT DO NOTHING")
		if _, err = query.RunWith(tx).Exec(); err != nil {
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
