package client

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	pb "github.com/kolya59/easy_normalization/proto"
)

func SendCars(cars []pb.Car, host, port string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect")
	}
	defer conn.Close()
	c := pb.NewCarSaverClient(conn)

	// Convert cars
	convertedCars := make([]*pb.Car, len(cars))
	for i, car := range cars {
		convertedCars[i] = &car
	}

	// Save cars
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SaveCars(ctx, &pb.SaveRequest{Cars: convertedCars})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to save cars")
	}
	log.Printf("Result: %s", r.GetMessage())
}
