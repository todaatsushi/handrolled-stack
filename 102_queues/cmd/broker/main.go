package broker

import (
	"log"

	"github.com/todaatsushi/queue/internal/broker"
)

func Run(port int) {
	log.Println("Starting broker")

	server := broker.NewServer(port)
	log.Fatal(server.Start())
}
