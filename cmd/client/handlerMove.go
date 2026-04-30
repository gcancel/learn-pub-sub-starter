package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerMove(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.ArmyMove) pubsub.Acktype {

	return func(am gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")

		moveOutcome := gs.HandleMove(am)
		switch moveOutcome {
		case gamelogic.MoveOutComeSafe:
			fmt.Println("message acknowledged.")
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:

			warQueue := routing.WarRecognitionsPrefix + "." + gs.GetUsername()
			err := pubsub.PublishJSON(ch, routing.ExchangePerilTopic, warQueue,
				gamelogic.RecognitionOfWar{
					Attacker: am.Player,
					Defender: gs.GetPlayerSnap(),
				})
			if err != nil {
				return pubsub.NackRequeue
			}

			fmt.Println("message acknowledged.")
			return pubsub.Ack
		case gamelogic.MoveOutcomeSamePlayer:
			fmt.Println("message discarded.")
			return pubsub.NackDiscard
		}
		fmt.Println("error: unknown move outcome")
		fmt.Println("message discarded.")
		return pubsub.NackDiscard
	}
}
