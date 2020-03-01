package client

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"

	pb "github.com/kolya59/easy_normalization/proto"
)

func SendCars(cars []pb.Car, url, topic string) error {
	if topic == "" {
		topic = "cars"
	}
	connection, err := amqp.Dial(url)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to RabbitMQ broker")
		return fmt.Errorf("failed to connect to RabbitMQ broker: %v", err)
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get channel")
		return fmt.Errorf("failed to get channel: %v", err)
	}

	if err := channel.ExchangeDeclare(topic, "topic", true, false, false, false, nil); err != nil {
		log.Error().Err(err).Msg("Failed to declare exchange channel")
		return fmt.Errorf("failed to declare exchange channel: %v", err)
	}

	data, err := json.Marshal(cars)
	if err != nil {
		log.Error().Err(err).Msg("Failed to publish data")
		return fmt.Errorf("failed to publish data: %v", err)
	}
	msg := amqp.Publishing{Body: data}

	if err := channel.Publish(topic, "random-key", false, false, msg); err != nil {
		log.Error().Err(err).Msg("Failed to publish data")
		return fmt.Errorf("failed to publish data: %v", err)
	}
	return nil
}
