package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/todaatsushi/basic_tcp/cmd/client"
	"github.com/todaatsushi/basic_tcp/cmd/server"
	"github.com/todaatsushi/basic_tcp/internal/encoding"
)

func main() {
	// Config
	runType := flag.String("type", "", "'CLIENT' or 'SERVER'")

	// Server
	port := flag.Int("p", 1337, "Connection port.")

	// Client
	msg := flag.String("m", "", "Message to server.")

	flag.Parse()

	if len(*runType) == 0 {
		log.Fatal("Please provide a type ('CLIENT' or 'SERVER').")
	}
	parsedType := strings.ToLower(*runType)

	const basePrefix = "handrolled::TCP"
	logPrefix := fmt.Sprintf("%s::%s", basePrefix, strings.ToLower(parsedType))
	log.SetPrefix(logPrefix + "\t")

	switch strings.ToLower(parsedType) {
	case "client":
		client.Send(*msg, *port, encoding.Basic{})
	case "server":
		server.Run(port, encoding.Basic{})
	default:
		log.Fatalf("Not a valid 'type' value (%s). Must be 'CLIENT' or 'SERVER.'", *runType)
	}
}
