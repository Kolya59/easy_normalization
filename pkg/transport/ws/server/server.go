package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"

	postgresdriver "github.com/kolya59/easy_normalization/pkg/postgres-driver"
	pb "github.com/kolya59/easy_normalization/proto"
)

var upgrader = websocket.Upgrader{} // use default options

func handler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade")
		return
	}
	defer c.Close()

	// Read data
	var cars []pb.Car
	var newCar pb.Car
	for err = c.ReadJSON(&newCar); err == nil; err = c.ReadJSON(&newCar) {
		cars = append(cars, newCar)
	}

	// Send data in DB
	if err = postgresdriver.SaveCars(cars); err != nil {
		log.Error().Err(err).Msg("Could not send cars to DB")
		return
	}
	log.Info().Msgf("Cars %v was saved via WS", cars)
}

func StartServer(host, port string, done chan interface{}) {
	r := chi.NewRouter()
	r.HandleFunc("/", handler)

	// TODO: Open WS server
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
	log.Info().Msgf("WS server is listening on %s:%s", host, port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("Failed to listen and serve")
	}
}
