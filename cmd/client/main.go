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
	mqttclient "github.com/kolya59/easy_normalization/pkg/transport/mqtt/client"
	restclient "github.com/kolya59/easy_normalization/pkg/transport/rest/client"
	wsclient "github.com/kolya59/easy_normalization/pkg/transport/ws/client"
)

var opts struct {
	Host          string `long:"host" env:"HOST" description:"Server host" required:"true"`
	Port          string `long:"port" env:"PORT" description:"Server port" required:"true"`
	RedisServer   string `long:"redis_server" env:"REDIS_SERVER" description:"Redis servers" required:"true"`
	RedisPassword string `long:"redis_password" env:"REDIS_PASSWORD" description:"Password for servers" required:"true"`
	RedisDatabase int    `long:"redis_database" env:"REDIS_DATABASE" description:"Redis database" required:"true"`
	BrokerHost    string `long:"host" env:"HOST" description:"Host" required:"true"`
	BrokerPort    int    `long:"port" env:"PORT" description:"Port" required:"true"`
	User          string `long:"user" env:"USER" description:"Username" required:"true"`
	Password      string `long:"password" env:"PASS" description:"Password" required:"true"`
	Topic         string `long:"topic" env:"TOPIC" description:"Topic" required:"true"`
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

	// Send data to server
	restclient.SendCars(car.Data, opts.Host, opts.Port)
	wsclient.SendCars(car.Data, opts.Host, opts.Port)
	mqttclient.SendCars(car.Data, opts.BrokerHost, opts.BrokerPort, opts.User, opts.Password, opts.Topic)
}
