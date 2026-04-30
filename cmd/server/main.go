package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const url = "amqp://guest:guest@localhost:5672/"

	fmt.Println("Starting Peril server...")
	connection, err := amqp.Dial(url)
	if err != nil {
		log.Fatal("Error connecting to the server. ", err)
	}
	defer connection.Close()

	fmt.Printf("Connection was successful!\n")

	ch, err := connection.Channel()
	if err != nil {
		log.Fatal(err)
	}

	err = pubsub.SubscribeGob(connection, routing.ExchangePerilTopic, routing.GameLogSlug, "game_logs.*", pubsub.Durable, handlerLogs())
	if err != nil {
		log.Fatal(err)
	}

	// printing Server Help
	gamelogic.PrintServerHelp()
	for {
		args := gamelogic.GetInput()
		if len(args) == 0 {
			continue
		}
		cmd := args[0]
		switch cmd {
		case "pause":
			fmt.Printf("Pausing the game...\n")

			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: true,
			})
		case "resume":
			fmt.Printf("Resuming game...\n")

			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: false,
			})
		case "quit":
			fmt.Printf("Exiting game")
			os.Exit(0)
		default:
			fmt.Printf("Cannot comprehend the command...")
		}

	}

	// waiting for os signals
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// <-signalChan

}
