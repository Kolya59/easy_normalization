package client

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"

	pb "github.com/kolya59/easy_normalization/proto"
)

const (
	RabbitMQUrl = "amqp://%s:%s@%s:%s/%s"
)

func SendCars(cars []pb.Car, brokerHost, brokerPort, user, password, topic string) {
	if topic == "" {
		topic = "cars"
	}
	connection, err := amqp.Dial(fmt.Sprintf(RabbitMQUrl, user, password, brokerHost, brokerPort, user))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to RabbitMQ broker")
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get channel")
	}

	if err := channel.ExchangeDeclare(topic, "topic", true, false, false, false, nil); err != nil {
		log.Fatal().Err(err).Msg("Failed to declare exchange channel")
	}

	data, err := json.Marshal(cars)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to publish data")
	}
	msg := amqp.Publishing{Body: data}

	if err := channel.Publish(topic, "random-key", false, false, msg); err != nil {
		log.Fatal().Err(err).Msg("Failed to publish data")
	}
}
