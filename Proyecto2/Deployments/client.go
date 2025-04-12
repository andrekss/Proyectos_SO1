package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	pb "github.com/andres/Proyecto2/proto" // importa tu paquete generado

	"google.golang.org/grpc"
)

type tweet struct {
	Description string `json:"Description"`
	Country     string `json:"Country"`
	Weather     string `json:"Weather"`
}

func sendToGrpc(tweet tweet) {
	conn, err := grpc.Dial("grpc-server:50051", grpc.WithInsecure()) // usa DNS del servicio Kubernetes
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTweetPublisherClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err = client.PublishToKafka(ctx, &pb.Tweet{
		Description: tweet.Description,
		Country:     tweet.Country,
		Weather:     tweet.Weather,
	})
	if err != nil {
		log.Printf("Error al enviar a Kafka: %v", err)
	}

	_, err = client.PublishToRabbit(ctx, &pb.Tweet{
		Description: tweet.Description,
		Country:     tweet.Country,
		Weather:     tweet.Weather,
	})
	if err != nil {
		log.Printf("Error al enviar a RabbitMQ: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	var tweet tweet
	err := json.NewDecoder(r.Body).Decode(&tweet)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	go sendToGrpc(tweet)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tweet recibido"))
}

func main() {
	http.HandleFunc("/tweet", handler)
	log.Println("API REST escuchando en :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
