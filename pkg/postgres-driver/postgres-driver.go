package postgresdriver

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	"github.com/psu/easy_normalization/pkg/car"
	"gitlab.com/gs-iot/gs-firmware/pkg/firmware"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

const (
	selectFirmwareQuery = "SELECT major_firmware_version, minor_firmware_version, bootloader_version, status, path FROM firmwares WHERE id=(SELECT max(firmwares.id) FROM firmwares WHERE bootloader_version = $1 AND major_firmware_version = $2)"
	insertFirmwareQuery = "INSERT INTO firmwares(major_firmware_version, minor_firmware_version, bootloader_version, status, path) VALUES ($1, $2, $3, $4, $5)"
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

// Select last firmware with defined version
func SelectFirmware(bootLoaderVersion int, firmwareVersion int) (sfw firmware.SerializableFirmware, err error) {
	selectFirmware, err := db.Prepare(selectFirmwareQuery)
	if err != nil {
		return firmware.SerializableFirmware{}, fmt.Errorf("could not prepare select query: %v", err)
	}
	defer func() {
		err = selectFirmware.Close()
		if err != nil {
			log.Error().Err(err).Msgf("Could not close database connection:")
		}
	}()
	err = selectFirmware.QueryRow(bootLoaderVersion, firmwareVersion).Scan(
		&sfw.MajorFirmwareVersion,
		&sfw.MinorFirmwareVersion,
		&sfw.BootloaderVersion,
		&sfw.Status,
		&sfw.Path,
	)
	if err != nil {
		return firmware.SerializableFirmware{}, fmt.Errorf("could not select firmware: %v", err)
	}
	return sfw, nil
}

// Insert new firmware into database
func InsertFirmware(bv int, majorFv int, minorFv int, status string, path string) (err error) {
	insertFirmware, err := db.Prepare(insertFirmwareQuery)
	if err != nil {
		return fmt.Errorf("could not prepare insert query: %v", err)
	}
	defer func() {
		err = insertFirmware.Close()
		if err != nil {
			log.Error().Err(err).Msgf("Could not close database connection:")
		}
	}()
	_, err = insertFirmware.Exec(majorFv, minorFv, bv, status, path)
	if err != nil {
		return fmt.Errorf("could not insert firmware into database: %v", err)
	}
	log.Info().Msgf("Firmware bv = %v majorFv = %v is added in database", bv, majorFv)
	return nil
}

func SendData(newCar *car.Car) error {
	return nil
}
