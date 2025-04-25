package main

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

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
			log.Printf("Mensaje recibido de RabbitMQ y enviado a valkey: %s", d.Body)
		}

	}

}

func main() {

	go ConsumerRabbitMq()
	select {}
}
