package broker

import (
	"log"
	"net"
)

type Server struct {
	port int
}

func NewServer(port int) *Server {
	return &Server{
		port: port,
	}
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

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
}
