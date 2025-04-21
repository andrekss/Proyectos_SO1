package main

import (
	"context"
	"log"
	"net"

	pb "github.com/andres/Proyecto2/proto" // Ajusta esto si tu módulo tiene otro path
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTweetPublisherServer
}

func (s *server) PublishToKafka(ctx context.Context, tweet *pb.Tweet) (*pb.Response, error) {
	log.Printf("Kafka <- Tweet: %+v", tweet)
	return &pb.Response{Status: "Publicado en Kafka"}, nil
}

func (s *server) PublishToRabbit(ctx context.Context, tweet *pb.Tweet) (*pb.Response, error) {
	log.Printf("RabbitMQ <- Tweet: %+v", tweet)
	return &pb.Response{Status: "Publicado en RabbitMQ"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("Error al escuchar en puerto 8081: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTweetPublisherServer(s, &server{})

	log.Println("gRPC server escuchando en puerto 8081...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
