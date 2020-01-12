package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/kolya59/easy_normalization/pkg/transport/mqtt/common"
	pb "github.com/kolya59/easy_normalization/proto"
)

func sendToBroker(uri *url.URL, topic string, cars []pb.Car) {
	client := common.Connect("sub", uri)
	defer client.Disconnect(uint(10 * time.Second.Seconds()))
	data, err := json.Marshal(cars)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to publish data")
	}
	client.Publish(topic, 0, false, data)
}

func SendCars(cars []pb.Car, brokerHost, brokerPort, user, password, topic string) {
	uri, err := url.Parse(fmt.Sprintf(common.CloudMQTTUrl, user, password, brokerHost, brokerPort, topic))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse uri")
	}
	if topic == "" {
		topic = "cars"
	}

	sendToBroker(uri, topic, cars)
}
