package main

import (
	"log"

	"github.com/streadway/amqp"
)

func main() {
	dialResult := amqpDial()
	defer do(dialResult, func(conn dial) {
		conn.Close() 
	}, nil)
	
	channelResult := fmap(dialResult, amqpChannel)
	defer do(channelResult, func(ch channel) {
		ch.Close() 
	}, nil)
	
	queueResult := fmap(channelResult, amqpQeue)

	consumeResult := fmap(channelResult, func(ch channel) result[consume] {
		return amqpConsume(ch)(queueResult)
	})

	success := func(cs consume) {
		log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
		for d := range cs.msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}
	do(consumeResult, success, exitOnFailure)
}

// type aliases required to make library types to work properly with generics as of now
type (
	dial struct {
		*amqp.Connection
	}
	channel struct {
		*amqp.Channel
	}
	queue struct {
		amqp.Queue
	}
	consume struct {
		msgs <-chan amqp.Delivery
	}
)

type result[T any] struct {
	value T
	err   error
}

func try[T any](fn func() (value T, err error)) result[T] {
	val, err := fn()
	return result[T]{value: val, err: err}
}

func fmap[T, B any](res result[T], fn func(T) result[B]) result[B] {
	if res.err != nil {
		return result[B]{value: zero[B](), err: res.err}
	}
	return fn(res.value)
}

func do[T any](res result[T], fn func(t T), errFn func(err error)) {
	if res.err != nil && errFn != nil {
		errFn(res.err)
		return
	}
	if fn != nil {
		fn(res.value)
	}
}

func zero[T any]() T {
	var z T
	return z
}

func exitOnFailure(err error) {
	log.Fatal(err)
}

func amqpDial() result[dial] {
	return try(func() (dial, error) {
		con, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		return dial{con}, err
	})
}

func amqpChannel(conn dial) result[channel] {
	return try(func() (channel, error) {
		ch, err := conn.Channel()
		return channel{ch}, err
	})
}

func amqpQeue(ch channel) result[queue] {
	return try(func() (queue, error) {
		q, err := ch.QueueDeclare(
			"hello", // name
			false,   // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		return queue{q}, err
	})
}

func amqpConsume(ch channel) func(result[queue]) result[consume] {
	return func(queueResult result[queue]) result[consume] {
		return fmap(queueResult, func(q queue) result[consume] {
			return try(func() (consume, error) {
				msgs, err := ch.Consume(
					q.Name, // queue
					"",     // consumer
					true,   // auto-ack
					false,  // exclusive
					false,  // no-local
					false,  // no-wait
					nil,    // args
				)
				return consume{msgs}, err
			})
		})
	}
}
