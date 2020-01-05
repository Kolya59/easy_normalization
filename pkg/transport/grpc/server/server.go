package server

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	postgresdriver "github.com/kolya59/easy_normalization/pkg/postgres-driver"
	pb "github.com/kolya59/easy_normalization/proto"
)

type server struct {
	pb.UnimplementedCarSaverServer
}

func (s *server) SaveCars(ctx context.Context, in *pb.SaveRequest) (*pb.SaveReply, error) {
	// Convert data
	cars := make([]pb.Car, len(in.Cars))
	for i, car := range in.Cars {
		cars[i] = *car
	}

	// Send data in DB
	if err := postgresdriver.SaveCars(cars); err != nil {
		log.Error().Err(err).Msg("Could not send cars to DB")
		return &pb.SaveReply{}, err
	}

	return &pb.SaveReply{Message: "All is ok"}, nil
}

func StartServer(host, port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}
	s := grpc.NewServer()
	pb.RegisterCarSaverServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("Failed to serve")
	}
}
