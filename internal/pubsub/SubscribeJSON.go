package pubsub

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype int

const (
	Ack Acktype = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T) Acktype) error {

	channel, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		log.Fatal(err)
	}
	consumeChan, err := channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for message := range consumeChan {
			var body T
			err := json.Unmarshal(message.Body, &body)
			if err != nil {
				log.Fatal(err)
			}
			switch handler(body) {
			case Ack:
				message.Ack(false)
				fmt.Println("Ack")
			case NackDiscard:
				fmt.Println("Nack")
				message.Nack(false, false)
			case NackRequeue:
				fmt.Println("Requeue")
				message.Nack(false, true)
			}

		}

	}()

	return nil
}
