package main

import (
	"context"
	"log"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
)

type Tweet struct {
	Description string
	Country     string
	Weather     string
}

var rdb *redis.Client

func ConectarRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Error al conectar a Redis: %v", err)
	}
}

func GuardarTweetRedis(tweet Tweet) {
	key := "tweet:" + tweet.Description
	err := rdb.HSet(context.Background(), key, map[string]interface{}{
		"description": tweet.Description,
		"country":     tweet.Country,
		"weather":     tweet.Weather,
	}).Err()

	if err != nil {
		log.Printf("Error al guardar tweet en Redis: %v", err)
	} else {
		log.Printf("Tweet guardado en Redis: %+v", tweet)
	}
}

func parsearTweet(mensaje string) (Tweet, error) {
	partes := strings.Split(mensaje, " - ")
	if len(partes) != 3 {
		return Tweet{}, log.Output(1, "Formato inválido del mensaje recibido")
	}
	return Tweet{
		Description: strings.TrimSpace(partes[0]),
		Country:     strings.TrimSpace(partes[1]),
		Weather:     strings.TrimSpace(partes[2]),
	}, nil
}

func ConsumerKafka() {
	ConectarRedis()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"kafka:9092"},
		Topic:     "tweets",
		GroupID:   "test-group",
		Partition: 0,
		MinBytes:  1,
		MaxBytes:  10e6,
	})

	log.Println("Kafka consumer escuchando...")

	for {
		//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		m, err := reader.ReadMessage(context.Background())
		//cancel()

		if err != nil {
			log.Printf("Error al leer mensaje: %v", err)
			continue
		}
		log.Printf("Mensaje recibido de Kafka: %s", string(m.Value))

		// Proceso de guardado
		tweet, err := parsearTweet(string(m.Value))
		if err != nil {
			log.Printf("Error al parsear tweet: %v", err)
			continue
		}

		GuardarTweetRedis(tweet)
	}
}

func main() {
	go ConsumerKafka()
	select {}
}
