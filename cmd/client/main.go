package main

import (
	_ "database/sql"
	"os"
	"time"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack"

	"github.com/kolya59/easy_normalization/pkg/car"
	postgresdriver "github.com/kolya59/easy_normalization/pkg/postgres-driver"
)

var opts struct {
	RedisServer   string `long:"redis_server" env:"REDIS_SERVER" description:"Redis servers" required:"true"`
	RedisPassword string `long:"redis_password" env:"REDIS_PASSWORD" description:"Password for servers" required:"true"`
	RedisDatabase int    `long:"redis_database" env:"REDIS_DATABASE" description:"Redis database" required:"true"`
	ProfilerPort  string `long:"prof_port" env:"PROF_PORT" description:"Profiler port" required:"false"`
	LogLevel      string `long:"log_level" env:"LOG_LEVEL" description:"Log level for zerolog" required:"false"`
}

// Send info to Redis database
func SetCar(newCar *car.Car, codec *cache.Codec, index string) error {
	err := codec.Set(&cache.Item{
		Ctx:        nil,
		Key:        index,
		Object:     newCar,
		Func:       nil,
		Expiration: time.Minute,
	})

	if err != nil {
		return err
	}
	return nil
}

// Get info from Redis data base
func GetCar(index string, codec *cache.Codec) (newCar *car.Car, err error) {
	newCar = &car.Car{}
	err = codec.Get(index, newCar)
	if err != nil {
		return nil, err
	}
	return newCar, nil
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

	// Redis initialization
	servers := map[string]string{
		"server1": "localhost:6379",
	}
	ring := redis.NewRing(&redis.RingOptions{
		Addrs:    servers,
		Password: "",
		DB:       0,
	})
	defer func() {
		err := ring.Close()
		if err != nil {
			log.Fatal().Msgf("Could not close ring: %v", err)
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
		err = SetCar(&obj, codec, string(i))
		if err != nil {
			log.Fatal().Msgf("Could not add car in Redis: %v", err)
		}
	}

	// Set DB structure
	err = postgresdriver.InitDatabaseStructure()
	if err != nil {
		log.Fatal().Msgf("Could not init Postgres structure: %v", err)
	}

	// Set data in DB
	for i, _ := range car.Data {
		// Get data fromm Redis
		newCar, err := GetCar(string(i), codec)
		if err != nil {
			log.Fatal().Msgf("Could not get car from Redis: %v", err)
		}
		// Send data in DB
		err = postgresdriver.SendData(newCar)
		if err != nil {
			log.Fatal().Msgf("Could not send car in Postgres: %v", err)
		}
	}
}
