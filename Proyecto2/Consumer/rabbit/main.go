package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
)

var ctx = context.Background()
var rdb *redis.Client

type Tweet struct {
	Description string
	Country     string
	Weather     string
}

func ConectarValkey() {
	for {
		rdb = redis.NewClient(&redis.Options{
			Addr: "valkey:6379", // Conectamos a Valkey
		})

		_, err := rdb.Ping(ctx).Result()
		if err != nil {
			log.Printf("Error al conectar a Valkey: %v. Reintentando en 1 segundos...", err)
			time.Sleep(1 * time.Second)
			continue
		}

		log.Println("Conectado exitosamente a Valkey")
		break
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

func GuardarTweetValkey(tweet Tweet) {

	// Incrementar contadores
	_, err := rdb.Incr(ctx, "messages:counter").Result()
	if err != nil {
		log.Printf("Error incrementando messages:counter en Valkey: %v", err)
	}

	_, err = rdb.HIncrBy(ctx, "country:counter", tweet.Country, 1).Result()
	if err != nil {
		log.Printf("Error incrementando country:counter en Valkey: %v", err)
	}

	// Guardar tweet en hash con ID único
	tweetID, err := rdb.Incr(ctx, "tweet:counter").Result()
	if err != nil {
		log.Printf("Error incrementando tweet:counter en Valkey: %v", err)
		return
	}

	key := fmt.Sprintf("tweet:%d", tweetID)

	err = rdb.HSet(ctx, key, map[string]interface{}{
		"description": tweet.Description,
		"country":     tweet.Country,
		"weather":     tweet.Weather,
	}).Err()

	if err != nil {
		log.Printf("Error guardando tweet en Valkey: %v", err)
	} else {
		log.Printf("Tweet guardado en Valkey:%+v bajo clave %s", tweet, key)
	}
}

func ConsumerRabbitMq() {
	var conn *amqp.Connection
	var err error

	for {
		conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
			log.Printf("Error al conectar con RabbitMQ: %v. Reintentando en 1 segundos...", err)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error al abrir canal: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"tweets",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error al declarar cola: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error al consumir cola: %v", err)
	}

	log.Println("RabbitMQ consumer escuchando...")
	for d := range msgs {
		msg := string(d.Body)
		tweet, err := parsearTweet(msg)
		if err != nil {
			log.Printf("Error al parsear tweet: %v", err)
			continue
		}

		log.Printf("Mensaje recibido de RabbitMQ: %s", msg)
		GuardarTweetValkey(tweet)
	}
}

func main() {
	ConectarValkey()
	go ConsumerRabbitMq()
	select {}
}
