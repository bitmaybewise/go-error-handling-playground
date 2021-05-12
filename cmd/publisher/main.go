package main

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// https://www.rabbitmq.com/tutorials/tutorial-one-go.html
func main() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	body := "Hello World!"
	for d := range time.Tick(time.Second * 10) {
		msg := fmt.Sprintf("%s - %s", d, body)
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			})
		if err != nil {
			log.Fatalf("%s: %s", "Failed to publish a message", err)
		}
		log.Printf(" [x] Sent %s", msg)
	}
}
