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
	grpcserver "github.com/kolya59/easy_normalization/pkg/transport/grpc/server"
	"github.com/kolya59/easy_normalization/pkg/transport/proxy"
	rabbitmqserver "github.com/kolya59/easy_normalization/pkg/transport/rabbitmq/server"
	restserver "github.com/kolya59/easy_normalization/pkg/transport/rest/server"
	wsserver "github.com/kolya59/easy_normalization/pkg/transport/ws/server"
)

var opts struct {
	Host       string `long:"host" env:"HOST" description:"Server host" required:"true"`
	RESTPort   string `long:"rest_port" env:"REST_PORT" description:"Server port" required:"true"`
	WSPort     string `long:"ws_port" env:"WS_PORT" description:"Server port" required:"true"`
	GRPCPort   string `long:"grpc_port" env:"GRPC_PORT" description:"Server port" required:"true"`
	DbHost     string `long:"database_host" env:"DB_HOST" description:"Database host" required:"true"`
	DbPort     string `long:"database_port" env:"DB_PORT" description:"Database port" required:"true"`
	DbName     string `long:"database_name" env:"DB_NAME" description:"Database name" required:"true"`
	DbUser     string `long:"database_username" env:"DB_USER" description:"Database username" required:"true"`
	DbPassword string `long:"database_password" env:"DB_PASSWORD" description:"Database password" required:"true"`
	LogLevel   string `long:"log_level" env:"LOG_LEVEL" description:"Log level for zerolog" required:"false"`
	BrokerHost string `long:"broker_host" env:"BROKER_HOST" description:"Host" required:"true"`
	BrokerPort string `long:"broker_port" env:"BROKER_PORT" description:"Port" required:"true"`
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
	if _, err := flags.ParseArgs(&opts, os.Args); err != nil {
		log.Panic().Msgf("Could not parse flags: %v", err)
	}

	level, err := zerolog.ParseLevel(opts.LogLevel)
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Connect to database
	log.Debug().Msg("Try to connect to database")
	if err = postgresdriver.InitDatabaseConnection(opts.DbHost, opts.DbPort, opts.DbUser, opts.DbPassword, opts.DbName); err != nil {
		log.Fatal().Msgf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err = postgresdriver.CloseConnection(); err != nil {
			log.Fatal().Msgf("Could not close db connection: %v", err)
		}
	}()

	// Graceful shutdown
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM)
	signal.Notify(sigint, syscall.SIGINT)
	done := make(chan interface{})

	// Start servers
	go proxy.StartServer()
	log.Info().Msg("Started Proxy server")
	go restserver.StartServer(opts.Host, opts.RESTPort, done)
	log.Info().Msg("Started REST server")
	go wsserver.StartServer(opts.Host, opts.WSPort, done)
	log.Info().Msg("Started WS server")
	go rabbitmqserver.StartServer(opts.BrokerHost, opts.BrokerPort, opts.User, opts.Password, opts.Topic, done)
	log.Info().Msg("Started RabbitMQ server")
	go grpcserver.StartServer(opts.Host, opts.GRPCPort)
	log.Info().Msg("Started GRPC server")

	// Wait interrupt signal
	select {
	case <-sigint:
		close(done)
	}
}
