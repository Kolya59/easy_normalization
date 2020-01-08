package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	if resp.StatusCode != http.StatusOK {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal().Msgf("Response status is not OK: %v\nFailed to read body", resp.Status)
		}
		log.Fatal().Msgf("Response status is not OK: %v\nBody: %s", resp.Status, respBody)
	}
}
