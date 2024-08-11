package server

import (
	"io"
	"net"
	"time"

	"github.com/todaatsushi/handrolled-cache/internal/cache"
)

type Server struct {
	store *cache.Store
}

func (s *Server) Run(port int) error {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: port})
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn io.ReadWriter) {
}

type c struct{}

func (clock c) Now() time.Time {
	return time.Now()
}

func (clock c) CalcExpires(ttl int) time.Time {
	return clock.Now().Add(time.Second * time.Duration(ttl))
}

func NewServer(cacheSize int) *Server {
	return &Server{store: cache.NewStore(uint64(cacheSize), c{})}
}
