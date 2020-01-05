package main

import (
	_ "database/sql"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	postgresdriver "github.com/kolya59/easy_normalization/pkg/postgres-driver"
	mqttserver "github.com/kolya59/easy_normalization/pkg/transport/mqtt/server"
	restserver "github.com/kolya59/easy_normalization/pkg/transport/rest/server"
	wsserver "github.com/kolya59/easy_normalization/pkg/transport/ws/server"
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
	BrokerHost string `long:"host" env:"HOST" description:"Host" required:"true"`
	BrokerPort string `long:"port" env:"PORT" description:"Port" required:"true"`
	User       string `long:"user" env:"USER" description:"Username" required:"true"`
	Password   string `long:"password" env:"PASS" description:"Password" required:"true"`
	Topic      string `long:"topic" env:"TOPIC" description:"Topic" required:"true"`
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

	// Graceful shutdown
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM)
	signal.Notify(sigint, syscall.SIGINT)
	done := make(chan interface{})

	// Start servers
	go restserver.StartServer(opts.Host, opts.Port, done)
	go wsserver.StartServer(opts.Host, opts.Port, done)
	go mqttserver.StartServer(opts.BrokerHost, opts.BrokerPort, opts.User, opts.Password, opts.Topic, done)

	<-sigint
	close(done)
}