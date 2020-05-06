package client

import (
	"context"
	_ "database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	grpcclient "github.com/kolya59/easy_normalization/pkg/transport/grpc/client"
	pubsubclient "github.com/kolya59/easy_normalization/pkg/transport/mq/client"
	restclient "github.com/kolya59/easy_normalization/pkg/transport/rest/client"
	wsclient "github.com/kolya59/easy_normalization/pkg/transport/ws/client"
	pb "github.com/kolya59/easy_normalization/proto"
)

var opts struct {
	Server    string `long:"server" env:"SERVER" required:"true"`
	Port      string `long:"port" env:"PORT" description:"Server port" required:"true"`
	ProjectID string `long:"projectID" env:"PROJECT_ID" required:"true" default:"trrp-virus"`
	RESTPort  string `long:"rest_port" env:"REST_PORT" description:"Server port" required:"true"`
	WSPort    string `long:"ws_port" env:"WS_PORT" description:"Server port" required:"true"`
	GRPCPort  string `long:"grpc_port" env:"GRPC_PORT" description:"Server port" required:"true"`
	LogLevel  string `long:"log_level" env:"LOG_LEVEL" description:"Log level for zerolog" required:"false"`
	Topic     string `long:"topic" env:"TOPIC" description:"Topic" required:"true"`
}

var (
	defaultCars []pb.Car
	client      *pubsubclient.Client
)

func fillData() []pb.Car {
	return []pb.Car{
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
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(200)
		_, _ = w.Write([]byte("Client is ready"))
		return
	case http.MethodPost:
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte("Failed to read data"))
			return
		}

		var cars []pb.Car
		if err = json.Unmarshal(data, &cars); err != nil {
			log.Error().Err(err).Msg("Failed to decode body")
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte("Failed to decode body"))
			return
		}

		t := r.Header.Get("Type")
		switch t {
		case "REST":
			if err := restclient.SendCars(cars, opts.Server, opts.RESTPort); err != nil {
				log.Error().Err(err).Msg("Failed to send cars via REST")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Failed to send cars via REST"))
				return
			}
		case "WS":
			if err := wsclient.SendCars(cars, opts.Server, opts.WSPort); err != nil {
				log.Error().Err(err).Msg("Failed to send cars via WS")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Failed to send cars via WS"))
				return
			}
		case "MQ":
			if err := client.SendCars(cars); err != nil {
				log.Error().Err(err).Msg("Failed to send cars via MQ")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Failed to send cars via MQ"))
				return
			}
		case "gRPC":
			if err := grpcclient.SendCars(cars, opts.Server, opts.GRPCPort); err != nil {
				log.Error().Err(err).Msg("Failed to send cars via gRPC")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Failed to send cars via gRPC"))
				return
			}
		case "All":
			if err := restclient.SendCars(defaultCars[:2], opts.Server, opts.RESTPort); err != nil {
				log.Error().Err(err).Msg("Failed to send cars via REST")
				return
			}
			if err := wsclient.SendCars(defaultCars[1:3], opts.Server, opts.WSPort); err != nil {
				log.Error().Err(err).Msg("Failed to send cars via WS")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Failed to send cars via WS"))
				return
			}
			if err := client.SendCars(defaultCars[2:4]); err != nil {
				log.Error().Err(err).Msg("Failed to send cars via AMPQ")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Failed to send cars via AMPQ"))
				return
			}
			if err := grpcclient.SendCars(defaultCars[3:], opts.Server, opts.GRPCPort); err != nil {
				log.Error().Err(err).Msg("Failed to send cars via gRPC")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Failed to send cars via gRPC"))
				return
			}
		default:
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte("Unrecognized header"))
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Method is not allowed"))
		return
	}

	log.Info().Msg("Cars sent successfully")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Cars sent successfully"))
}

func Start(done chan interface{}) {
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
		log.Fatal().Err(err).Msg("Could not parse flags")
	}

	level, err := zerolog.ParseLevel(opts.LogLevel)
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	defaultCars = fillData()

	if client, err = pubsubclient.NewClient(opts.ProjectID, opts.Topic); err != nil {
		log.Fatal().Err(err).Msg("Failed to create new client")
	}

	r := http.NewServeMux()

	r.HandleFunc("/", handler)

	srv := http.Server{
		Addr:         fmt.Sprintf(":%s", opts.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		ctx := context.Background()
		<-done
		_ = srv.Shutdown(ctx)
	}()

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("Failed to listen and serve")
	}
}
