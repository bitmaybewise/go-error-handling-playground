package main

import (
	"log"

	performer "github.com/hlmerscher/go-error-handling-playground/perfomer"
)

// https://www.rabbitmq.com/tutorials/tutorial-one-go.html
func main() {
	dial := &performer.Dial{}
	channel := &performer.Channel{Dial: dial}
	queue := &performer.QueueDeclare{Ch: channel}
	consume := &performer.Consume{Ch: channel, QueueDeclare: queue}
	err := performer.Do(
		dial.Do,
		channel.Do,
		queue.Do,
		consume.Do,
	)
	defer dial.Conn.Close()
	defer channel.Ch.Close()

	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range consume.Messages {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
