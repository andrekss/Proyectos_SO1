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

func Obtenemos_Tweet_Rust() *tweet {
	var tweet *tweet

	tRust, err := getTweetFromRust()
	if err != nil {
		log.Printf("Error al obtener tweet desde Rust: %v", err)
		time.Sleep(2 * time.Second) // espera no enciclar
	} else {
		tweet = tRust
		log.Printf(" Tweet obtenido desde Rust: %+v", tweet)
		// sendToGrpc(t)
		//break
	}

	return tweet
}

// Comunicación gRPC

func sendToGrpcKafka(tweet tweet) {

	// ----- Conectando con grpc -----
	conn, err := grpc.Dial("kafka-writer:8083", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTweetPublisherClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// ----- Conectando con grpc -----

	// ----- Publicación kafka -----
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
	// ----- Publicación kafka -----

}

func sendToGrpcRabbit(tweet tweet) {

	// ----- Conectando con grpc -----
	conn, err := grpc.Dial("rabbit-writer:8084", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTweetPublisherClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// ----- Conectando con grpc -----

	// ----- Publicación rabbit -----
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
	// ----- Publicación rabbit -----
}

func sendTweetToKafka(tweet_recibido *tweet) {

	go sendToGrpcKafka(*tweet_recibido)

	log.Printf("Tweet enviado a Kafka: %s -- %s -- %s", tweet_recibido.Description, tweet_recibido.Country, tweet_recibido.Weather)
}

func sendTweetToRabbit(tweet_recibido *tweet) {

	go sendToGrpcRabbit(*tweet_recibido)

	log.Printf("Tweet enviado a Rabbit: %s -- %s -- %s", tweet_recibido.Description, tweet_recibido.Country, tweet_recibido.Weather)
}

func main() {

	var tweet_recibido *tweet

	for {
		tweet_recibido = Obtenemos_Tweet_Rust()

		if tweet_recibido != nil {
			log.Printf("Tweet recibido: %+v", tweet_recibido)
			break
		}

		log.Println("Esperando tweet desde Rust...")
		time.Sleep(2 * time.Second)
	}

	// Publicamos a los brokers
	go sendTweetToKafka(tweet_recibido)
	go sendTweetToRabbit(tweet_recibido)

	log.Println("API REST escuchando en :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
