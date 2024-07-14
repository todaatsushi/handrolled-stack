package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

func main() {
	runType := flag.String("type", "", "'CLIENT' or 'SERVER'")
	flag.Parse()

	if len(*runType) == 0 {
		log.Fatal("Please provide a type ('CLIENT' or 'SERVER').")
	}
	parsedType := strings.ToLower(*runType)

	const basePrefix = "handrolled::TCP"
	logPrefix := fmt.Sprintf("%s::%s", basePrefix, strings.ToLower(parsedType))
	log.SetPrefix(logPrefix)

	switch strings.ToLower(parsedType) {
	case "client":
		log.Fatal("TODO")
	case "server":
		log.Fatal("TODO")
	default:
		log.Fatalf("Not a valid 'type' value (%s). Must be 'CLIENT' or 'SERVER.'", *runType)
	}
}
