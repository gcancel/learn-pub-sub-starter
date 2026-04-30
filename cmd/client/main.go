package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	url := "amqp://guest:guest@localhost:5672/"

	fmt.Println("Starting Peril client...")

	connection, err := amqp.Dial(url)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := connection.Channel()
	if err != nil {
		log.Fatal(err)
	}

	message, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal(err)
	}

	// creating a new game state
	gameState := gamelogic.NewGameState(message)

	//_, _, err = pubsub.DeclareAndBind(connection, routing.ExchangePerilDirect, username, routing.PauseKey, pubsub.Transient)
	err = pubsub.SubscribeJSON(connection, routing.ExchangePerilDirect, routing.PauseKey+"."+gameState.GetUsername(), routing.PauseKey, pubsub.Transient, handlerPause(gameState))
	if err != nil {
		log.Fatal(err)
	}

	err = pubsub.SubscribeJSON(connection, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+gameState.GetUsername(), "army_moves.*", pubsub.Transient, handlerMove(gameState, ch))
	if err != nil {
		log.Fatal(err)
	}

	// consuming all war messages:
	err = pubsub.SubscribeJSON(connection, routing.ExchangePerilTopic, routing.WarRecognitionsPrefix, "war.*", pubsub.Durable, handlerWar(gameState, ch))
	if err != nil {
		log.Fatal(err)
	}

	for {

		args := gamelogic.GetInput()
		if len(args) == 0 {
			continue
		}

		cmd := args[0]
		switch cmd {
		case "spawn":
			err := gameState.CommandSpawn(args)
			if err != nil {
				log.Fatal(err)
			}
		case "move":
			move, err := gameState.CommandMove(args)
			if err != nil {
				log.Fatal(err)
			}
			//player := gameState.GetPlayerSnap()
			//units := gameState.GetUnitsSnap()
			//location := args[1]
			err = pubsub.PublishJSON(ch, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+gameState.GetUsername(), move)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Move was successful!", move)

		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			if args[1] == "" {
				fmt.Println("value need to spam")
				continue
			}
			num, err := strconv.Atoi(args[1])
			if err != nil {
				log.Fatal(err)
			}
			for range num {
				msg := gamelogic.GetMaliciousLog()
				err := pubsub.PublishGob(ch, routing.ExchangePerilTopic, routing.GameLogSlug+"."+gameState.GetUsername(), msg)
				if err != nil {
					log.Fatal(err)
				}
			}

		case "quit":
			fmt.Println("Exiting game.")
			os.Exit(0)
		default:
			fmt.Println("Invalid commands, please retry")
		}
	}

	// waiting for os signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

}

func checkWords(str []string) bool {
	validLocations := map[string]bool{"americas": true, "europe": true, "africa": true, "asia": true, "antarctica": true, "australia": true}
	validUnits := map[string]bool{"infantry": true, "cavalry": true, "artillery": true}

	location := str[0]
	unit := str[1]

	_, ok := validLocations[location]
	if !ok {
		return false
	}
	_, ok = validUnits[unit]
	if !ok {
		return false
	}
	return true

}
