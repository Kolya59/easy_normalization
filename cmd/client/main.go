package main

import (
	_ "database/sql"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack"

	grpcclient "github.com/kolya59/easy_normalization/pkg/transport/grpc/client"
	rabbitmqclient "github.com/kolya59/easy_normalization/pkg/transport/rabbitmq/client"
	restclient "github.com/kolya59/easy_normalization/pkg/transport/rest/client"
	wsclient "github.com/kolya59/easy_normalization/pkg/transport/ws/client"
	pb "github.com/kolya59/easy_normalization/proto"
)

var opts struct {
	Host       string `long:"host" env:"HOST" description:"Server host" required:"true"`
	RESTPort   string `long:"rest_port" env:"REST_PORT" description:"Server port" required:"true"`
	WSPort     string `long:"ws_port" env:"WS_PORT" description:"Server port" required:"true"`
	GRPCPort   string `long:"grpc_port" env:"GRPC_PORT" description:"Server port" required:"true"`
	LogLevel   string `long:"log_level" env:"LOG_LEVEL" description:"Log level for zerolog" required:"false"`
	BrokerHost string `long:"broker_host" env:"BROKER_HOST" description:"Host" required:"true"`
	BrokerPort string `long:"broker_port" env:"BROKER_PORT" description:"Port" required:"true"`
	User       string `long:"user" env:"USER" description:"Username" required:"true"`
	Password   string `long:"password" env:"PASS" description:"Password" required:"true"`
	Topic      string `long:"topic" env:"TOPIC" description:"Topic" required:"true"`
}

// Send info to Redis database
func SetCar(newCar *pb.Car, codec *cache.Codec, index string) error {
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
func GetCar(index string, codec *cache.Codec) (newCar *pb.Car, err error) {
	newCar = &pb.Car{}
	err = codec.Get(index, newCar)
	if err != nil {
		return nil, err
	}
	return newCar, nil
}

func fillData() []pb.Car {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": ":6379",
		},
	})

	codec := &cache.Codec{
		Redis: ring,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}

	cars := []pb.Car{
		{
			Model:                   "2114",
			BrandName:               "LADA",
			BrandCreatorCountry:     "Russia",
			EngineModel:             "V123",
			EnginePower:             80,
			EngineVolume:            16,
			EngineType:              "L4",
			TransmissionModel:       "M123",
			TransmissionType:        "M",
			TransmissionGearsNumber: 5,
			WheelModel:              "Luchshie kolesa Rossii",
			WheelRadius:             13,
			WheelColor:              "Black",
			Price:                   120000,
		},
		{
			Model:                   "2115",
			BrandName:               "LADA",
			BrandCreatorCountry:     "Russia",
			EngineModel:             "V124",
			EnginePower:             100,
			EngineVolume:            18,
			EngineType:              "L4",
			TransmissionModel:       "M123",
			TransmissionType:        "M",
			TransmissionGearsNumber: 5,
			WheelModel:              "Luchshie kolesa Rossii",
			WheelRadius:             13,
			WheelColor:              "Black",
			Price:                   150000,
		},
		{
			Model:                   "Rio",
			BrandName:               "Kia",
			BrandCreatorCountry:     "Korea",
			EngineModel:             "V14234",
			EnginePower:             100,
			EngineVolume:            90,
			EngineType:              "V4",
			TransmissionModel:       "A123",
			TransmissionType:        "A",
			TransmissionGearsNumber: 4,
			WheelModel:              "Luchie kolesa Kitaya",
			WheelRadius:             15,
			WheelColor:              "Red",
			Price:                   400000,
		},
		{
			Model:                   "Sportage",
			BrandName:               "Kia",
			BrandCreatorCountry:     "Korea",
			EngineModel:             "V14234",
			EnginePower:             100,
			EngineVolume:            90,
			EngineType:              "V4",
			TransmissionModel:       "A1234",
			TransmissionType:        "A",
			TransmissionGearsNumber: 5,
			WheelModel:              "Luchie kolesa Kitaya",
			WheelRadius:             15,
			WheelColor:              "Red",
			Price:                   400000,
		},
		{
			Model:                   "A500",
			BrandName:               "Mercedes",
			BrandCreatorCountry:     "Germany",
			EngineModel:             "E1488",
			EnginePower:             300,
			EngineVolume:            50,
			EngineType:              "V12",
			TransmissionModel:       "R123",
			TransmissionType:        "A",
			TransmissionGearsNumber: 8,
			WheelModel:              "Luchshie kolesa Armenii",
			WheelRadius:             20,
			WheelColor:              "Green",
			Price:                   3000000,
		},
	}
	for i, car := range cars {
		if err := SetCar(&car, codec, fmt.Sprintf("%v", i)); err != nil {
			log.Fatal().Err(err).Msg("Failed to set car")
		}
	}

	res := make([]pb.Car, len(cars))
	for i := 0; i < len(cars); i++ {
		car, err := GetCar(fmt.Sprintf("%v", i), codec)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get car")
		}
		res[i] = *car
	}

	return res
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

	cars := fillData()

	// Send data to server
	restclient.SendCars(cars[:2], opts.Host, opts.RESTPort)
	wsclient.SendCars(cars[1:3], opts.Host, opts.WSPort)
	rabbitmqclient.SendCars(cars[2:4], opts.BrokerHost, opts.BrokerPort, opts.User, opts.Password, opts.Topic)
	grpcclient.SendCars(cars[3:], opts.Host, opts.GRPCPort)
}
