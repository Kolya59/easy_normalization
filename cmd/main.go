package main

import (
	_ "database/sql"
	"os"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/jessevdk/go-flags"
	"github.com/psu/easy_normalization/pkg/car"
	postgresdriver "github.com/psu/easy_normalization/pkg/postgres-driver"
	redisdriver "github.com/psu/easy_normalization/pkg/redis-driver"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack"
)

var opts struct {
	DbHost        string `long:"database_host" env:"DB_HOST" description:"Database host" required:"true"`
	DbPort        string `long:"database_port" env:"DB_PORT" description:"Database port" required:"true"`
	DbName        string `long:"database_name" env:"DB_NAME" description:"Database name" required:"true"`
	DbUser        string `long:"database_username" env:"DB_USER" description:"Database username" required:"true"`
	DbPassword    string `long:"database_password" env:"DB_PASSWORD" description:"Database password" required:"true"`
	RedisServer   string `long:"redis_server" env:"REDIS_SERVER" description:"Redis servers" required:"true"`
	RedisPassword string `long:"redis_password" env:"REDIS_PASSWORD" description:"Password for servers" required:"true"`
	RedisDatabase int    `long:"redis_database" env:"REDIS_DATABASE" description:"Redis database" required:"true"`
	ProfilerPort  string `long:"prof_port" env:"PROF_PORT" description:"Profiler port" required:"false"`
	LogLevel      string `long:"log_level" env:"LOG_LEVEL" description:"Log level for zerolog" required:"false"`
}

func main() {
	// Log initialization
	zerolog.MessageFieldName = "MESSAGE"
	zerolog.LevelFieldName = "LEVEL"
	zerolog.ErrorFieldName = "ERROR"
	zerolog.TimestampFieldName = "TIME"
	zerolog.CallerFieldName = "CALLER"
	log.Logger = log.Output(os.Stderr).With().Str("PROGRAM", "firmware-update-server").Caller().Logger()

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
			log.Error().Msgf("Could not close db connection: %v", err)
		}
	}()

	// Redis initialization
	servers := map[string]string{
		"server1": opts.RedisServer,
	}
	ring := redis.NewRing(&redis.RingOptions{
		Addrs:    servers,
		Password: opts.RedisPassword,
		DB:       opts.RedisDatabase,
	})
	defer func() {
		err = ring.Close()
		if err != nil {
			log.Error().Msgf("Could not close ring: %v", err)
		}
	}()
	codec := &cache.Codec{
		Redis:     ring,
		Marshal:   msgpack.Marshal,
		Unmarshal: msgpack.Unmarshal,
	}

	// Set data in Redis
	car.FillData()
	for i, obj := range car.Data {
		err = redisdriver.SetCar(&obj, codec, string(i))
		if err != nil {
			log.Error().Msgf("Could not add car in Redis: %v", err)
		}
	}

	// Set DB structure
	postgresdriver.InitDatabaseStructure()

	for i, _ := range car.Data {
		// Get data fromm Redis
		newCar, err := redisdriver.GetCar(codec, string(i))
		if err != nil {
			log.Error().Msgf("Could not get car from Redis: %v", err)
		}
		// Send data in DB
		err = postgresdriver.SendData(newCar)
		if err != nil {
			log.Error().Msgf("Could not send car in Postgres: %v", err)
		}
	}
}
