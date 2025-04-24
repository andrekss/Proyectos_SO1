package consumer_kafka

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"

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
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		m, err := reader.ReadMessage(ctx)
		cancel()

		if err != nil {
			log.Printf("Error al leer mensaje: %v", err)
			continue
		}

		log.Printf("Mensaje recibido de Kafka: %s", string(m.Value))
	}
}

func PublicarEnKafka() string {
	resp, err := http.Get("http://0.0.0.0:8081/publishKafka")
	if err != nil {
		log.Fatalf("Error al hacer la petición HTTP: %v", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error al leer la respuesta HTTP: %v", err)
		return ""
	}

	return string(body)
}

func main() {
	tweet := PublicarEnKafka()
	print(tweet)
	go ConsumerKafka()
	select {}
}
