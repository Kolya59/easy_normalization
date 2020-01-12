package server

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"

	postgresdriver "github.com/kolya59/easy_normalization/pkg/postgres-driver"
	pb "github.com/kolya59/easy_normalization/proto"
)

const (
	RabbitMQUrl = "amqp://%s:%s@%s:%s/%s"
)

func handleConnection(msg amqp.Delivery) {
	// Get data
	var cars []pb.Car
	if err := json.Unmarshal(msg.Body, cars); err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal cars")
	}
	// Send data in DB
	if err := postgresdriver.SaveCars(cars); err != nil {
		log.Error().Err(err).Msg("Could not send cars to DB")
	}
	log.Info().Msgf("Cars %v was saved via RabbitMQ", cars)
	msg.Ack(false)
}

func StartServer(brokerHost, brokerPort, user, password, topic string, done chan interface{}) {
	if topic == "" {
		topic = "cars"
	}
	connection, err := amqp.Dial(fmt.Sprintf(RabbitMQUrl, user, password, brokerHost, brokerPort, user))
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to RabbitMQ broker")
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get channel")
	}

	if _, err = channel.QueueDeclare("cars", true, false, false, false, nil); err != nil {
		log.Error().Err(err).Msg("Failed to declare the queue")
	}

	if err = channel.QueueBind("cars", "#", topic, false, nil); err != nil {
		panic("error binding to the queue: " + err.Error())
	}

	msgs, err := channel.Consume("cars", "", false, false, false, false, nil)
	if err != nil {
		panic("error consuming the queue: " + err.Error())
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
