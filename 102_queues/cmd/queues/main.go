package main

import (
	"flag"
	"log"
	"strings"

	"github.com/todaatsushi/queue/cmd/broker"
	"github.com/todaatsushi/queue/cmd/producer"
)

func main() {
	port := flag.Int("port", 1337, "Port")
	messages := flag.String("m", "", "Comma separated messages to send as tasks.")
	runType := flag.String("type", "", "'BROKER' | 'PRODUCER' | 'HEALTH'")
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
	default:
		log.Printf("Unhandled run type '%s'", *runType)
	}
}
