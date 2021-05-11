package performer

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Task func() error

func Do(it ...Task) error {
	for _, fn := range it {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

type Dial struct {
	Conn *amqp.Connection
}

func (m *Dial) Do() (err error) {
	log.Println("Dial.Do")
	m.Conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	return
}

type Channel struct {
	Dial *Dial
	Ch   *amqp.Channel
}

func (m *Channel) Do() (err error) {
	log.Println("Channel.Do")
	m.Ch, err = m.Dial.Conn.Channel()
	return
}

type QueueDeclare struct {
	Ch    *Channel
	Queue amqp.Queue
}

func (m *QueueDeclare) Do() (err error) {
	log.Println("QueueDeclare.Do")
	m.Queue, err = m.Ch.Ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	return
}

type Consume struct {
	Ch           *Channel
	QueueDeclare *QueueDeclare
	Messages     <-chan amqp.Delivery
}

func (m *Consume) Do() (err error) {
	log.Println("Consume.Do")
	m.Messages, err = m.Ch.Ch.Consume(
		m.QueueDeclare.Queue.Name, // queue
		"",                        // consumer
		true,                      // auto-ack
		false,                     // exclusive
		false,                     // no-local
		false,                     // no-wait
		nil,                       // args
	)
	return
}

type Publish struct {
	Ch           *Channel
	QueueDeclare *QueueDeclare
}

func (m *Publish) Do() (err error) {
	log.Println("Publish.Do")
	body := "Hello World!"
	for v := range time.Tick(time.Second * 10) {
		msg := fmt.Sprintf("%s - %s", body, v)
		log.Println("sending message ->", msg)
		err = m.Ch.Ch.Publish(
			"",                        // exchange
			m.QueueDeclare.Queue.Name, // routing key
			false,                     // mandatory
			false,                     // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			})
		if err != nil {
			log.Fatalf("%s: %s", "Failed to publish a message", err)
		}
	}
	log.Printf(" [x] Sent %s", body)
	return
}
