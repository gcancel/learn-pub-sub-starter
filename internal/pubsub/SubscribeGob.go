package pubsub

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeGob[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T) Acktype) error {
	channel, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		log.Fatal(err)
	}
	channel.Qos(10, 0, true)
	consumeChan, err := channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for message := range consumeChan {
			var body T
			data := bytes.NewBuffer(message.Body)
			dec := gob.NewDecoder(data)
			dec.Decode(&body)

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
