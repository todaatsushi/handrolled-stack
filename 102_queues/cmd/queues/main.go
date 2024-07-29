package main

import (
	"flag"

	"github.com/todaatsushi/queue/cmd/broker"
)

func main() {
	port := flag.Int("port", 1337, "Port")
	broker.Run(*port)
}
