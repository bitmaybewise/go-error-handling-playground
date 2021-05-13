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
	publish := tasks.Publish{Ch: channel, QueueDeclare: queue}
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
