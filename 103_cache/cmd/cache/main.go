package main

import (
	"flag"
	"log"
	"strings"

	"github.com/todaatsushi/handrolled-cache/cmd/client"
	"github.com/todaatsushi/handrolled-cache/cmd/server"
)

func main() {
	port := flag.Int("p", 420, "Runs on.")
	cacheSize := flag.Int("c", 0, "Max number of items in cache.")
	runType := flag.String("type", "", "One of 'SERVER' or 'CLIENT'")

	flag.Parse()

	t := strings.ToLower(*runType)

	switch t {
	case "server":
		log.Fatal(server.Start(*port, *cacheSize))
	case "client":
		log.Fatal(client.Start(*port))
	default:
		log.Fatal("'type' must be one of 'SERVER' or 'CLIENT'.")
	}
}
