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

var ctx = context.Background() // necesario para redis v8

func GuardarTweetValkey(tweet string) {
	parts := strings.Split(tweet, " - ")
	if len(parts) != 3 {
		log.Println("Formato incorrecto del tweet para Valkey")
		return
	}
	description := strings.TrimSpace(parts[0])
	country := strings.TrimSpace(parts[1])
	weather := strings.TrimSpace(parts[2])

	rdb := redis.NewClient(&redis.Options{
		Addr: "valkey:6379",
	})

	key := fmt.Sprintf("tweet:%d", time.Now().UnixNano())

	err := rdb.HSet(ctx, key, map[string]interface{}{
		"description": description,
		"country":     country,
		"weather":     weather,
	}).Err()

	if err != nil {
		log.Printf("Error guardando en Valkey: %v", err)
	} else {
		log.Printf("Guardado en Valkey bajo clave %s", key)
	}
}

func ConsumerRabbitMq() {
	var conn *amqp.Connection
	var err error
	for {
		conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
			log.Printf("Error al conectar con RabbitMQ: %v. Reintentando en 2 segundos...", err)
			time.Sleep(2 * time.Second)
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
		"tweets", // mismo nombre que se usa al publicar
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error al declarar cola: %v", err)
	}

	for {
		msgs, err := ch.Consume(
			q.Name,
			"",
			true, // auto-ack
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
			log.Printf("Mensaje recibido de RabbitMQ: %s", d.Body)
			GuardarTweetValkey(msg)
		}

	}

}

func main() {

	go ConsumerRabbitMq()
	select {}
}
