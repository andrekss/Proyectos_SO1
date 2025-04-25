package main

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func ConsumerKafka() {
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

		log.Printf("Mensaje recibido de Kafka y enviado a redis: %s", string(m.Value))
	}
}

func main() {
	go ConsumerKafka()
	select {}
}
