package server

import (
	"encoding/json"
	"fmt"
	"net/url"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"

	postgresdriver "github.com/kolya59/easy_normalization/pkg/postgres-driver"
	"github.com/kolya59/easy_normalization/pkg/transport/mqtt/common"
	pb "github.com/kolya59/easy_normalization/proto"
)

func listen(uri *url.URL, topic string) {
	client := common.Connect("sub", uri)
	client.Subscribe(topic, 0, handleConnection)

	// TODO: Set deadline
	/*select {
	case <-time.NewTimer(1000 * time.Second).C:
		client.Unsubscribe(topic)
		return
	}*/
}

func handleConnection(client mqtt.Client, msg mqtt.Message) {
	// Get data
	var cars []pb.Car
	if err := json.Unmarshal(msg.Payload(), cars); err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal cars")
	}
	// Send data in DB
	if err := postgresdriver.SaveCars(cars); err != nil {
		log.Error().Err(err).Msg("Could not send cars to DB")
	}
}

func StartServer(brokerHost, brokerPort, user, password, topic string, done chan interface{}) {
	uri, err := url.Parse(fmt.Sprintf(common.CloudMQTTUrl, user, password, brokerHost, brokerPort, topic))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse url")
	}
	if topic == "" {
		topic = "test"
	}

	listen(uri, topic)
}
