package server

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"

	postgresdriver "github.com/kolya59/easy_normalization/pkg/postgres-driver"
	pb "github.com/kolya59/easy_normalization/proto"
)

func handleConnection(msg amqp.Delivery) {
	// Get data
	var cars []pb.Car
	if err := json.Unmarshal(msg.Body, &cars); err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal cars")
	}
	// Send data in DB
	if err := postgresdriver.SaveCars(cars); err != nil {
		log.Error().Err(err).Msg("Could not send cars to DB")
	}
	log.Info().Msgf("Cars %v was saved via RabbitMQ", cars)
	msg.Ack(false)
}

func StartServer(url, topic string, done chan interface{}) {
	if topic == "" {
		topic = "cars"
	}
	connection, err := amqp.Dial(url)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to RabbitMQ broker")
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get channel")
	}

	if err := channel.ExchangeDeclare(topic, "fanout", true, false, false, false, nil); err != nil {
		log.Fatal().Err(err).Msg("Failed to declare exchange channel")
	}

	if _, err = channel.QueueDeclare("", true, false, false, false, nil); err != nil {
		log.Error().Err(err).Msg("Failed to declare the queue")
	}

	if err = channel.QueueBind("", "#", topic, false, nil); err != nil {
		log.Error().Err(err).Msg("Failed to binding to the queue")
	}

	msgs, err := channel.Consume("", "", false, false, false, false, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to consuming the queue")
	}

loop:
	for {
		select {
		case <-done:
			break loop
		case msg := <-msgs:
			handleConnection(msg)
		}
	}
}
