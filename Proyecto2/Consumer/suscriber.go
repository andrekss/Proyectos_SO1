package consumer

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/streadway/amqp"
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

func ConsumerRabbitMq() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Error al conectar con RabbitMQ: %v", err)
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
		log.Printf("Mensaje recibido de RabbitMQ: %s", d.Body)
	}
}

func fetchTweetFromAPI() string {
	resp, err := http.Get("http://0.0.0.0:8081/tweet")
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

	tweet := fetchTweetFromAPI()

	print(tweet)
	go ConsumerKafka()
	go ConsumerRabbitMq()
	select {}
}
