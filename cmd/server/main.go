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
	DBURL           string `long:"db_url" env:"DATABASE_URL" description:"DB URL" required:"true"`
	CloudamqpUrl    string `long:"cloudamqp_url" env:"CLOUDAMQP_URL" description:"CLOUDAMQP URL" required:"true"`
	CloudamqpApikey string `long:"cloudamqp_apikey" env:"CLOUDAMQP_APIKEY" description:"CLOUDAMQP APIKEY" required:"true"`
	RESTPort        string `long:"rest_port" env:"PORT" description:"Server port" required:"true"`
	WSPort          string `long:"ws_port" env:"WS_PORT" description:"Server port" required:"true"`
	GRPCPort        string `long:"grpc_port" env:"GRPC_PORT" description:"Server port" required:"true"`
	LogLevel        string `long:"log_level" env:"LOG_LEVEL" description:"Log level for zerolog" required:"false"`
	Topic           string `long:"topic" env:"TOPIC" description:"Topic" required:"true"`
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
	if err = postgresdriver.InitDatabaseConnection(opts.DBURL); err != nil {
		log.Fatal().Msgf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err = postgresdriver.CloseConnection(); err != nil {
			log.Fatal().Msgf("Could not close db connection: %v", err)
		}
	}()

	/*// Set DB structure?
	log.Debug().Msg("Try to set database structure")
	if err = postgresdriver.InitDatabaseStructure(); err != nil {
		log.Fatal().Msgf("Could not init Postgres structure: %v", err)
	}*/

	// Graceful shutdown
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM)
	signal.Notify(sigint, syscall.SIGINT)
	done := make(chan interface{})

	// Start servers
	go proxy.StartServer()
	log.Info().Msg("Started Proxy server")
	go restserver.StartServer("", opts.RESTPort, done)
	log.Info().Msg("Started REST server")
	go wsserver.StartServer("", opts.WSPort, done)
	log.Info().Msg("Started WS server")
	// go mqttserver.StartServer(opts.BrokerHost, opts.BrokerPort, opts.User, opts.Password, opts.Topic, done)
	// log.Info().Msg("Started MQTT server")
	go rabbitmqserver.StartServer(opts.CloudamqpUrl, opts.Topic, done)
	log.Info().Msg("Started RabbitMQ server")
	go grpcserver.StartServer("", opts.GRPCPort)
	log.Info().Msg("Started GRPC server")

	// Wait interrupt signal
	select {
	case <-sigint:
		close(done)
	}
}
