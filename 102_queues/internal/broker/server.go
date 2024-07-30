package broker

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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
		port:  port,
		queue: make(chan messages.Message, 1000),
	}
}

func (s *Server) QueueLen() int {
	return len(s.queue)
}

func (s *Server) GetQueuedMessage() messages.Message {
	return <-s.queue
}

func (s *Server) Start() error {
	log.Printf("Starting server on port %v", s.port)

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: s.port})
	if err != nil {
		return err
	}
	log.Println("Starting listener.")
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		log.Println("New connection made.")

		go handle(conn, s)
	}
}

func (s *Server) ProcessMessage(w io.Writer, m messages.Message) error {
	switch m.Command {
	case messages.Log:
		log.Println("LOG:", m.Message)
	case messages.Enqueue:
		s.queue <- m
		log.Println("Message added to queue.")
	case messages.Consume:
		if len(m.Message) != 0 {
			return errors.New("Message should contain no data.")
		}
		if len(s.queue) == 0 {
			return nil
		}

		toConsume := s.GetQueuedMessage()
		data, err := toConsume.MarshalBinary()
		if err != nil {
			return err
		}

		_, err = w.Write(data)
		if err != nil {
			return err
		}
	case messages.QueueLen:
		numTasks := s.QueueLen()
		message := messages.NewMessage(messages.QueueLen, fmt.Sprint(numTasks))
		data, err := message.MarshalBinary()
		if err != nil {
			return err
		}
		_, err = w.Write(data)
		if err != nil {
			return err
		}

	default:
		panic("Unhandled")
	}
	return nil
}

func handle(conn net.Conn, server *Server) {
	log.Println("Handling connection.")

	reader := bufio.NewReader(conn)
	data, err := reader.ReadBytes(messages.DELIM)
	if err != nil {
		log.Println(err.Error())
		return
	}

	message, err := messages.UnmarshalBinary(data)
	if err != nil {
		log.Println(err.Error())
		return
	}
	server.ProcessMessage(conn, message)
}
