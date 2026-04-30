package pubsub

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Creating Ack type:

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {

	value, err := json.Marshal(val)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	mandatory := false
	immediate := false
	err = ch.PublishWithContext(ctx, exchange, key, mandatory, immediate, amqp.Publishing{
		ContentType: "application/json",
		Body:        value,
	})
	if err != nil {
		log.Fatal(err)
	}
	return nil

}
