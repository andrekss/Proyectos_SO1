package main

import (
	"context"
	"log"
	"net"

	pb "github.com/andres/Proyecto2/proto"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTweetPublisherServer
}

func (s *server) PublishToRabbit(ctx context.Context, tweet *pb.Tweet) (*pb.Response, error) {
	log.Printf("RabbitMQ <- Tweet: %+v", tweet)
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")

	if err != nil {
		log.Printf("Error al conectar con RabbitMQ: %v", err)
		return &pb.Response{Status: "Error RabbitMQ"}, err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return &pb.Response{Status: "Error RabbitMQ"}, err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"tweets", // nombre de la cola
		false,    // durable
		false,    // auto-delete
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return &pb.Response{Status: "Error RabbitMQ"}, err
	}

	body := tweet.Description + " - " + tweet.Country + " - " + tweet.Weather
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		return &pb.Response{Status: "Error RabbitMQ"}, err
	}

	return &pb.Response{Status: "Publicado en RabbitMQ"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8084")
	if err != nil {
		log.Fatalf("Error al escuchar en puerto 8084: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTweetPublisherServer(s, &server{})

	log.Println("gRPC server escuchando en puerto 8084...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
