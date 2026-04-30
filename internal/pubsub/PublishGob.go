package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	err := enc.Encode(val)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	mandatory := false
	immediate := false
	err = ch.PublishWithContext(ctx, exchange, key, mandatory, immediate, amqp.Publishing{
		ContentType: "application/gob",
		Body:        data.Bytes(),
	})
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
