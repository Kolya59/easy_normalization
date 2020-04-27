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

func SendCars(cars []pb.Car, host, port string) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%s", host, port), Path: "/"}
	log.Info().Msgf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to dial")
		return fmt.Errorf("failed to dial: %v", err)
	}
	defer c.Close()

	for _, newCar := range cars {
		select {
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			if err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
				log.Error().Err(err).Msg("Failed to write close msg")
				return fmt.Errorf("failed to write close msg: %v", err)
			}
			select {
			case <-time.After(time.Second):
			}
			return nil
		default:
			if err = c.WriteJSON(newCar); err != nil {
				log.Error().Err(err).Msg("Failed to write msg")
				return fmt.Errorf("failed to write msg: %v", err)
			}
		}
	}
	if err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		log.Error().Err(err).Msg("Failed to write close msg")
		return fmt.Errorf("failed to write close msg: %v", err)
	}

	return nil
}
