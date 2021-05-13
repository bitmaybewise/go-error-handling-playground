package main

import (
	"log"

	"github.com/hlmerscher/go-error-handling-playground/tasks"
	"github.com/hlmerscher/performer"
)

func main() {
	dial := &tasks.Dial{}
	channel := &tasks.Channel{Dial: dial}
	queue := &tasks.QueueDeclare{Ch: channel}
	consume := &tasks.Consume{Ch: channel, QueueDeclare: queue}
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

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	for d := range consume.Messages {
		log.Printf("Received a message: %s", d.Body)
	}
}
