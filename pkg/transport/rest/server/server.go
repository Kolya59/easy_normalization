package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"

	postgresdriver "github.com/kolya59/easy_normalization/pkg/postgres-driver"
	pb "github.com/kolya59/easy_normalization/proto"
)

func postCar(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get body")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Failed to read body"))
		return
	}
	defer r.Body.Close()
	var cars []pb.Car
	if err = json.Unmarshal(data, &cars); err != nil {
		log.Error().Err(err).Msg("Failed to decode body")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("Failed to decode body"))
		return
	}

	// Send data in DB
	if err = postgresdriver.SaveCars(cars); err != nil {
		log.Error().Err(err).Msg("Could not send cars to DB")
	}
	log.Info().Msgf("Cars %v was saved via REST", cars)
	w.WriteHeader(http.StatusOK)
}

func StartServer(host, port string, done chan interface{}) {
	// Set rest server
	r := chi.NewRouter()
	r.Post("/", postCar)

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: r,
		// TODO: TLS
		TLSConfig:    nil,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		// TODO Shutdown?
	}

	// Start server
	log.Info().Msgf("Server is listening on %s:%s", host, port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("Failed to listen and serve")
	}
}
