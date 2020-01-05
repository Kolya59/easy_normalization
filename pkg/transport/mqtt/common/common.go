package common

import (
	"fmt"
	"net/url"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
)

const (
	CloudMQTTUrl = "mqtt://%s:%s@%s:%s/%s"
)

func CreateClientOptions(clientId string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	return opts
}

func Connect(clientId string, uri *url.URL) mqtt.Client {
	opts := CreateClientOptions(clientId, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to broker")
	}
	return client
}
