package server

import (
	"errors"
	"flag"
	"log"
	"net"
)

func Run() error {
	log.Println("Starting server.")
	port := flag.Int("p", 1337, "Connection port.")
	flag.Parse()

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: *port})
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		_, err := listener.Accept()
		if err != nil {
			return err
		}
		return errors.New("TODO: implement binary parsing and log response.")
	}
}
