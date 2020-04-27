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
	queries := make([]sq.InsertBuilder, 5)
	queries[0] = sq.Insert("engines").
		Columns("engine_model", "engine_power", "engine_volume", "engine_type").
		Suffix("ON CONFLICT (engine_model) DO NOTHING")
	queries[1] = sq.Insert("transmissions").
		Columns("transmission_model", "transmission_type", "transmission_gears_number").
		Suffix("ON CONFLICT (transmission_model) DO NOTHING")
	queries[2] = sq.Insert("brands").
		Columns("brand_name", "brand_creator_country").
		Suffix("ON CONFLICT (brand_name) DO NOTHING")
	queries[3] = sq.Insert("wheels").
		Columns("wheel_model", "wheel_radius", "wheel_color").
		Suffix("ON CONFLICT (wheel_model) DO NOTHING")
	queries[4] = sq.Insert("cars").
		Columns("model", "engine", "transmission", "brand", "wheel", "price").
		Suffix("ON CONFLICT (model) DO NOTHING")

	// Bind arguments to queries
	for _, c := range cars {
		queries[0] = queries[0].
			Values(c.EngineModel, c.EnginePower, c.EngineVolume, c.EngineType)
		queries[1] = queries[1].
			Values(c.TransmissionModel, c.TransmissionType, c.TransmissionGearsNumber)
		queries[2] = queries[2].
			Values(c.BrandName, c.BrandCreatorCountry)
		queries[3] = queries[3].
			Values(c.WheelModel, c.WheelRadius, c.WheelColor)
		queries[4] = queries[4].
			Values(c.Model, c.EngineModel, c.TransmissionModel, c.BrandName, c.WheelModel, c.Price)
	}

	// Execute queries
	for _, query := range queries {
		q, a, err := query.PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			log.Error().Err(err).Msg("Failed to change placeholder")
			return err
		}
		if _, err = tx.Exec(q, a...); err != nil {
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
