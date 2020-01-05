package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"

	"github.com/kolya59/easy_normalization/pkg/car"
	postgresdriver "github.com/kolya59/easy_normalization/pkg/postgres-driver"
)

var upgrader = websocket.Upgrader{} // use default options

func handler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade")
		return
	}
	defer c.Close()
	var cars []car.Car

	// Read data
	// TODO: Set close handler
	for {
		var newCar car.Car
		if err := c.ReadJSON(&newCar); err != nil {
			log.Error().Err(err).Msg("Failed to read msg")
			break
		}
		cars = append(cars, newCar)
	}

	// Send data in DB
	if err = postgresdriver.SaveCars(cars); err != nil {
		log.Error().Err(err).Msg("Could not send cars to DB")
	}
}

func StartServer(host, port string, done chan interface{}) {
	r := chi.NewRouter()
	r.HandleFunc("/ws", handler)

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
