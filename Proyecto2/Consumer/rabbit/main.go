package consumer_rabbit

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/streadway/amqp"
)

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

func PublicarEnRabbit() string {
	resp, err := http.Get("http://0.0.0.0:8081/publishRabitt")
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

	tweet := PublicarEnRabbit()

	print(tweet)
	go ConsumerRabbitMq()
	select {}
}
