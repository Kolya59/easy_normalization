package client

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/rs/zerolog/log"

	"github.com/kolya59/easy_normalization/pkg/car"
	"github.com/kolya59/easy_normalization/pkg/transport/mqtt/common"
)

func sendToBroker(uri *url.URL, topic string, cars []car.Car) {
	client := common.Connect("sub", uri)
	for _, c := range cars {
		data, err := json.Marshal(c)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to publish data")
		}
		client.Publish(topic, 0, false, data)
	}
}

func SendCars(cars []car.Car, brokerHost, brokerPort, user, password, topic string) {
	uri, err := url.Parse(fmt.Sprintf(common.CloudMQTTUrl, user, password, brokerHost, brokerPort, topic))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse uri")
	}
	if topic == "" {
		topic = "time"
	}

	sendToBroker(uri, topic, cars)
}
