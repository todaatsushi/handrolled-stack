package broker

import (
	"log"

	"github.com/todaatsushi/queue/internal/broker"
)

func Run(port int) {
	server := broker.NewServer(port)
	log.Fatal(server.Start())
}
