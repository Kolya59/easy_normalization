package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"

	pb "github.com/kolya59/easy_normalization/proto"
)

func SendCars(cars []pb.Car, host, port string) {
	// Send data to server
	data, err := json.Marshal(cars)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to marshal data")
	}
	u := url.URL{Scheme: "http", Host: fmt.Sprintf("%s:%s", host, port), Path: "/"}
	resp, err := http.Post(u.String(), "application/json", bytes.NewReader(data))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to send request")
	}
	if resp.StatusCode == http.StatusOK {
		log.Info().Msgf("Response status is OK: %v", resp.Status)
	} else {
		log.Fatal().Msgf("Response status is not OK: %v", resp.Status)
	}
}
