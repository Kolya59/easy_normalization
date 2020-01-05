package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"

	postgresdriver "github.com/kolya59/easy_normalization/pkg/postgres-driver"
	pb "github.com/kolya59/easy_normalization/proto"
)

func postCar(w http.ResponseWriter, r *http.Request) {
	type PostReq struct {
		cars []pb.Car
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
