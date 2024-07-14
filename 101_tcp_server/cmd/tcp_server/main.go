package main

import (
	"log"

	"github.com/todaatsushi/handrolled/tcp_server/internal/server"
)

func main() {
	const logPrefix = "handrolled::TCP::server"
	log.SetPrefix(logPrefix + "\t")

	err := server.Run()
	if err != nil {
		panic(err)
	}
}
