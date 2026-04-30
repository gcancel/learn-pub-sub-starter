package pubsub

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	Durable SimpleQueueType = iota
	Transient
)

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType) (*amqp.Channel, amqp.Queue, error) {

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	queue, err := ch.QueueDeclare(
		queueName,
		isDurable(queueType),
		isTransient(queueType),
		isTransient(queueType),
		false,
		amqp.Table{
			"x-dead-letter-exchange": "peril_dlx",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	return ch, queue, nil
}

func isTransient(q SimpleQueueType) bool {
	if q == Transient {
		return true
	}
	return false
}

func isDurable(q SimpleQueueType) bool {
	if q == Durable {
		return true
	}
	return false
}
