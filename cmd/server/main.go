package main

import (
	_ "database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/kolya59/easy_normalization/pkg/car"
	postgresdriver "github.com/kolya59/easy_normalization/pkg/postgres-driver"
)

var opts struct {
	Host       string `long:"host" env:"HOST" description:"Server host" required:"true"`
	Port       string `long:"port" env:"PORT" description:"Server port" required:"true"`
	DbHost     string `long:"database_host" env:"DB_HOST" description:"Database host" required:"true"`
	DbPort     string `long:"database_port" env:"DB_PORT" description:"Database port" required:"true"`
	DbName     string `long:"database_name" env:"DB_NAME" description:"Database name" required:"true"`
	DbUser     string `long:"database_username" env:"DB_USER" description:"Database username" required:"true"`
	DbPassword string `long:"database_password" env:"DB_PASSWORD" description:"Database password" required:"true"`
	LogLevel   string `long:"log_level" env:"LOG_LEVEL" description:"Log level for zerolog" required:"false"`
}

func PostCar(w http.ResponseWriter, r *http.Request) {
	type PostReq struct {
		cars []car.Car
	}
	// Decode request
	req := &PostReq{}
	var data []byte
	_, err := r.Body.Read(data)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get body")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Failed to read body"))
		return
	}
	err = json.Unmarshal(data, req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode body")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("Failed to decode body"))
		return
	}

	// Send data in DB
	if err = postgresdriver.SaveCars(req.cars); err != nil {
		log.Error().Err(err).Msg("Could not send cars to DB")
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	// Log initialization
	zerolog.MessageFieldName = "MESSAGE"
	zerolog.LevelFieldName = "LEVEL"
	zerolog.ErrorFieldName = "ERROR"
	zerolog.TimestampFieldName = "TIME"
	zerolog.CallerFieldName = "CALLER"
	log.Logger = log.Output(os.Stderr).With().Str("PROGRAM", "easy-normalization").Caller().Logger()

	// Parse flags
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Panic().Msgf("Could not parse flags: %v", err)
	}

	level, err := zerolog.ParseLevel(opts.LogLevel)
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Connect to database
	err = postgresdriver.InitDatabaseConnection(opts.DbHost, opts.DbPort, opts.DbUser, opts.DbPassword, opts.DbName)
	if err != nil {
		log.Panic().Msgf("Failed to connect to database: %v", err)
	}
	defer func() {
		err = postgresdriver.CloseConnection()
		if err != nil {
			log.Fatal().Msgf("Could not close db connection: %v", err)
		}
	}()

	// TODO: Set DB structure?
	err = postgresdriver.InitDatabaseStructure()
	if err != nil {
		log.Fatal().Msgf("Could not init Postgres structure: %v", err)
	}

	// Set rest server
	r := chi.NewRouter()
	r.Post("/", PostCar)

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", opts.Host, opts.Port),
		Handler: r,
		// TODO: TLS
		TLSConfig:    nil,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Start server
	log.Info().Msgf("Server is listening on %s:%s", opts.Host, opts.Port)
	if err = srv.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("Failed to listen and serve")
	}
}
