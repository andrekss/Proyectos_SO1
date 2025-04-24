package main

import (
	"context"
	"log"
	"net"
	"strconv"

	pb "github.com/andres/Proyecto2/proto"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTweetPublisherServer
}

func (s *server) PublishToKafka(ctx context.Context, tweet *pb.Tweet) (*pb.Response, error) {
	log.Printf("Kafka <- Tweet: %+v", tweet)

	// Crear el topic si no existe
	conn, err := kafka.Dial("tcp", "kafka:9092")
	if err != nil {
		log.Printf("Error al conectar con Kafka: %v", err)
		return &pb.Response{Status: "Error Kafka Dial"}, err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Printf("Error obteniendo controlador Kafka: %v", err)
		return &pb.Response{Status: "Error Kafka Controller"}, err
	}

	controllerConn, err := kafka.Dial("tcp", controller.Host+":"+strconv.Itoa(controller.Port))
	if err != nil {
		log.Printf("Error conectando al broker controlador: %v", err)
		return &pb.Response{Status: "Error Kafka Broker Connect"}, err
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             "tweets",
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	if err != nil {
		log.Printf("Error al crear topic (o ya existe): %v", err)
	}

	// publicación a kafka
	writer := kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"), // nombre del servicio Kafka + puerto
		Topic:    "tweets",
		Balancer: &kafka.LeastBytes{},
	}

	msg := kafka.Message{
		Value: []byte(tweet.Description + " - " + tweet.Country + " - " + tweet.Weather),
	}

	err = writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf(" Error al enviar a Kafka: %v", err)
		return &pb.Response{Status: "Error Kafka"}, err
	}

	return &pb.Response{Status: "Publicado en Kafka"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Fatalf("Error al escuchar en puerto 8083: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTweetPublisherServer(s, &server{})

	log.Println("gRPC server escuchando en puerto 8083...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
