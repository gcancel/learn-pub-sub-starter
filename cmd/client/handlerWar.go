package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerWar(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.Acktype {
	return func(rw gamelogic.RecognitionOfWar) pubsub.Acktype {
		defer fmt.Print("> ")

		outcome, winner, loser := gs.HandleWar(rw)

		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			logs := fmt.Sprintf("%s won a war against %s", winner, loser)
			ack, err := publishGameLog(ch, gs, logs)
			if err != nil {
				log.Fatal(err)
			}
			return ack
		case gamelogic.WarOutcomeYouWon:
			logs := fmt.Sprintf("%s won a war against %s", winner, loser)
			ack, err := publishGameLog(ch, gs, logs)
			if err != nil {
				log.Fatal(err)
			}
			return ack
		case gamelogic.WarOutcomeDraw:
			logs := fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
			ack, err := publishGameLog(ch, gs, logs)
			if err != nil {
				log.Fatal(err)
			}
			return ack
		default:
			fmt.Println("error: outcome not defined.")
			return pubsub.NackDiscard
		}

	}

}

func publishGameLog(ch *amqp.Channel, gs *gamelogic.GameState, logs string) (pubsub.Acktype, error) {
	key := routing.GameLogSlug + "." + gs.GetUsername()
	err := pubsub.PublishGob(ch, routing.ExchangePerilTopic, key, routing.GameLog{
		CurrentTime: time.Now(),
		Message:     logs,
		Username:    gs.GetUsername(),
	})
	if err != nil {
		return pubsub.NackRequeue, err
	}
	return pubsub.Ack, nil
}
