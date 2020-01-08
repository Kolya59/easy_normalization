package client

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"

	pb "github.com/kolya59/easy_normalization/proto"
)

func SendCars(cars []pb.Car, host, port string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%s", host, port), Path: "/"}
	log.Info().Msgf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to dial")
	}
	defer c.Close()

	// TODO ????
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to read msg")
			}
			log.Info().Msgf("Recv: %s", message)
		}
	}()

	for _, newCar := range cars {
		select {
		case <-done:
			return
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to write close msg")
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		default:
			if err = c.WriteJSON(newCar); err != nil {
				log.Fatal().Err(err).Msg("Failed to write msg")
			}
		}
	}
}
