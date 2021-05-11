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
	publish := performer.Publish{Ch: channel, QueueDeclare: queue}
	err := performer.Do(
		dial.Do,
		channel.Do,
		queue.Do,
		publish.Do,
	)
	defer dial.Conn.Close()
	defer channel.Ch.Close()

	if err != nil {
		log.Fatal(err)
	}
}
