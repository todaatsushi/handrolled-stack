package main

import (
	"flag"
	"log"
	"strings"

	"github.com/todaatsushi/queue/cmd/broker"
	"github.com/todaatsushi/queue/cmd/consumer"
	"github.com/todaatsushi/queue/cmd/producer"
)

func main() {
	port := flag.Int("port", 1337, "Port")
	messages := flag.String("m", "", "Comma separated messages to send as tasks.")
	runType := flag.String("type", "", "'BROKER' | 'PRODUCER' | 'HEALTH' | 'QUEUELEN' | 'CONSUMER'")
	numConsumers := flag.Int("c", 1, "Num consumers")
	flag.Parse()

	_type := strings.ToUpper(*runType)
	switch _type {
	case "BROKER":
		log.Println("Starting broker.")
		broker.Run(*port)
	case "PRODUCER":
		log.Println("Sending messages.")
		producer.Send(*port, *messages)
	case "HEALTH":
		log.Println("Health check.")
		producer.Health(*port)
	case "QUEUELEN":
		log.Println("Getting length of queue.")
		producer.QueueLen(*port)
	case "CONSUMER":
		log.Printf("Starting %d consumers.\n", *numConsumers)
		consumer.Start(*port, *numConsumers)
	default:
		log.Printf("Unhandled run type '%s'", *runType)
	}
}
