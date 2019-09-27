package postgresdriver

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	"github.com/psu/easy_normalization/pkg/car"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

const (
	insertEngineQuery       = "INSERT INTO easy_normalization.normal.engines(engine_model, engine_power, engine_volume, engine_type) VALUES ($1, $2, $3, $4)"
	selectEngineQuery       = "SELECT engine_id FROM easy_normalization.normal.engines WHERE engine_model = $1"
	insertTransmissionQuery = "INSERT INTO easy_normalization.normal.transmissions(transmission_model, transmission_type, transmission_gears_number) VALUES ($1, $2, $3)"
	selectTransmissionQuery = "SELECT engine_id FROM easy_normalization.normal.engines WHERE engine_model = $1"
	insertBrandQuery        = "INSERT INTO easy_normalization.normal.brands(brand_name, brand_creator_country) VALUES ($1, $2)"
	selectBrandQuery        = "SELECT engine_id FROM easy_normalization.normal.engines WHERE engine_model = $1"
	insertWheelQuery        = "INSERT INTO easy_normalization.normal.wheels(wheels_model, wheels_radius, wheels_color) VALUES ($1, $2, $3)"
	selectWheelQuery        = "SELECT engine_id FROM easy_normalization.normal.wheels WHERE wheels_model = $1"
	insertCarQuery          = "INSERT INTO easy_normalization.normal.cars(model, brand_id, engine_id, transmission_id, wheel_id, price) VALUES ($1, $2, $3, $4, $5, $6)"
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
func SendData(newCar *car.Car) error {
	insertEngine, err := db.Prepare(insertEngineQuery)
	if err != nil {
		return fmt.Errorf("could not prepare insert query: %v", err)
	}
	defer func() {
		err = insertEngine.Close()
		if err != nil {
			log.Error().Err(err).Msgf("Could not close database connection:")
		}
	}()
	_, err = insertEngine.Exec(newCar.EngineModel, newCar.EnginePower, newCar.EngineVolume, newCar.EngineType)
	if err != nil {
		return fmt.Errorf("could not insert engine into database: %v", err)
	}
	selectEngine, err := db.Prepare(selectEngineQuery)
	if err != nil {
		return fmt.Errorf("could not prepare select query: %v", err)
	}
	defer func() {
		err = selectEngine.Close()
		if err != nil {
			log.Error().Err(err).Msgf("Could not close database connection:")
		}
	}()
	var engineId int
	err = selectEngine.QueryRow(newCar.EngineModel).Scan(&engineId)
	if err != nil {
		return fmt.Errorf("could not select firmware: %v", err)
	}
	log.Info().Msgf("Engine %v, %v, %v, %v is added in database", newCar.EngineModel, newCar.EnginePower, newCar.EngineVolume, newCar.EngineType)

	insertTransmission, err := db.Prepare(insertTransmissionQuery)
	if err != nil {
		return fmt.Errorf("could not prepare insert query: %v", err)
	}
	defer func() {
		err = insertTransmission.Close()
		if err != nil {
			log.Error().Err(err).Msgf("Could not close database connection:")
		}
	}()
	_, err = insertTransmission.Exec(newCar.TransmissionModel, newCar.TransmissionType, newCar.TransmissionGearsNumber)
	if err != nil {
		return fmt.Errorf("could not insert transmission into database: %v", err)
	}
	selectTransmission, err := db.Prepare(selectTransmissionQuery)
	if err != nil {
		return fmt.Errorf("could not prepare select query: %v", err)
	}
	defer func() {
		err = selectTransmission.Close()
		if err != nil {
			log.Error().Err(err).Msgf("Could not close database connection:")
		}
	}()
	var transmissionId int
	err = selectTransmission.QueryRow(newCar.TransmissionModel).Scan(&transmissionId)
	if err != nil {
		return fmt.Errorf("could not select firmware: %v", err)
	}
	log.Info().Msgf("Transmission %v, %v, %v is added in database", newCar.TransmissionModel, newCar.TransmissionType, newCar.TransmissionGearsNumber)

	insertBrand, err := db.Prepare(insertBrandQuery)
	if err != nil {
		return fmt.Errorf("could not prepare insert query: %v", err)
	}
	defer func() {
		err = insertBrand.Close()
		if err != nil {
			log.Error().Err(err).Msgf("Could not close database connection:")
		}
	}()
	_, err = insertBrand.Exec(newCar.BrandName, newCar.BrandCreatorCountry)
	if err != nil {
		return fmt.Errorf("could not insert brand into database: %v", err)
	}
	selectBrand, err := db.Prepare(selectBrandQuery)
	if err != nil {
		return fmt.Errorf("could not prepare select query: %v", err)
	}
	defer func() {
		err = selectBrand.Close()
		if err != nil {
			log.Error().Err(err).Msgf("Could not close database connection:")
		}
	}()
	var brandId int
	err = selectBrand.QueryRow(newCar.BrandName).Scan(&brandId)
	if err != nil {
		return fmt.Errorf("could not select firmware: %v", err)
	}
	log.Info().Msgf("Brand %v, %v is added in database", newCar.BrandName, newCar.BrandCreatorCountry)

	insertWheel, err := db.Prepare(insertWheelQuery)
	if err != nil {
		return fmt.Errorf("could not prepare insert query: %v", err)
	}
	defer func() {
		err = insertBrand.Close()
		if err != nil {
			log.Error().Err(err).Msgf("Could not close database connection:")
		}
	}()
	_, err = insertWheel.Exec(newCar.WheelModel, newCar.WheelRadius, newCar.WheelColor)
	if err != nil {
		return fmt.Errorf("could not insert wheel into database: %v", err)
	}
	selectWheel, err := db.Prepare(selectWheelQuery)
	if err != nil {
		return fmt.Errorf("could not prepare select query: %v", err)
	}
	defer func() {
		err = selectWheel.Close()
		if err != nil {
			log.Error().Err(err).Msgf("Could not close database connection:")
		}
	}()
	var wheelId int
	err = selectEngine.QueryRow(newCar.WheelModel).Scan(&wheelId)
	if err != nil {
		return fmt.Errorf("could not select firmware: %v", err)
	}
	log.Info().Msgf("Wheel %v, %v, %v is added in database", newCar.WheelModel, newCar.WheelRadius, newCar.WheelColor)

	insertCar, err := db.Prepare(insertCarQuery)
	if err != nil {
		return fmt.Errorf("could not prepare insert query: %v", err)
	}
	defer func() {
		err = insertCar.Close()
		if err != nil {
			log.Error().Err(err).Msgf("Could not close database connection:")
		}
	}()
	_, err = insertCar.Exec(newCar.Model, engineId, transmissionId, brandId, wheelId, newCar.Price)
	if err != nil {
		return fmt.Errorf("could not insert car into database: %v", err)
	}
	log.Info().Msgf("Car %v, %v, %v, %v, %v is added in database")
	return nil
}
