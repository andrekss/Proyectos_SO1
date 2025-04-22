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

// Comunicación con rust
func getTweetFromRust() (*tweet, error) {
	resp, err := http.Get("http://api-rust:8082/get_tweet") // nombre del contenedor Rust y puerto
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var t tweet
	err = json.NewDecoder(resp.Body).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Comunicación gRPC
func sendToGrpc(tweet tweet) {
	conn, err := grpc.Dial("grpc-server:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTweetPublisherClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	respPublishKafka, err := client.PublishToKafka(ctx, &pb.Tweet{
		Description: tweet.Description,
		Country:     tweet.Country,
		Weather:     tweet.Weather,
	})
	if err != nil {
		log.Printf("Error al enviar a Kafka: %v", err)
	} else {
		log.Printf("Respuesta Kafka: %v", respPublishKafka.GetStatus())
	}

	respPublishRabbit, err := client.PublishToRabbit(ctx, &pb.Tweet{
		Description: tweet.Description,
		Country:     tweet.Country,
		Weather:     tweet.Weather,
	})
	if err != nil {
		log.Printf("Error al enviar a RabbitMQ: %v", err)
	} else {
		log.Printf("Respuesta RabbitMQ: %v", respPublishRabbit.GetStatus())
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	var tweet_recibido *tweet
	for {
		tRust, err := getTweetFromRust()
		if err != nil {
			log.Printf("Error al obtener tweet desde Rust: %v", err)
			time.Sleep(2 * time.Second) // espera no enciclar
		} else {
			tweet_recibido = tRust
			log.Printf(" Tweet obtenido desde Rust: %+v", tweet_recibido)
			// sendToGrpc(t)
			break
		}
	}

	err := json.NewDecoder(r.Body).Decode(&tweet_recibido)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	go sendToGrpc(*tweet_recibido)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tweet recibido"))
}

func main() {
	http.HandleFunc("/tweet", handler) // Comunicación gRPC
	log.Println("API REST escuchando en :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
