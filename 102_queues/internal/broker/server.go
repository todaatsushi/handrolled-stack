package broker

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/todaatsushi/queue/internal/messages"
)

type Server struct {
	port  int
	queue chan messages.Message
}

func NewServer(port int) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) QueueLen() int {
	return len(s.queue)
}

func (s *Server) Start() error {
	log.Printf("Starting server on port %v", s.port)

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: s.port})
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go handle(conn, s)
	}
}

func (s *Server) ProcessMessage(m messages.Message) {
	log.Println("Starting")
	switch m.Command {
	case messages.Log:
		log.Println("LOG:", m.Message)
	case messages.Enqueue:
		go func() {
			s.queue <- m
		}()
	case messages.Consume:
		panic("TODO")
	}
}

func handle(conn net.Conn, server *Server) {
	reader := bufio.NewReader(conn)
	data, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Fprint(conn, err.Error())
		return
	}

	message, err := messages.UnmarshalBinary(data)
	if err != nil {
		fmt.Fprint(conn, err.Error())
		return
	}
	server.ProcessMessage(message)
}
