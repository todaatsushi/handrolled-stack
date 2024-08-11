package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/todaatsushi/handrolled-cache/internal/cache"
	"github.com/todaatsushi/handrolled-cache/internal/protocol"
)

type Server struct {
	store *cache.Store
}

func (s *Server) Run(port int) error {
	log.Println("Starting server on port", port)

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: port})
	if err != nil {
		return err
	}
	defer listener.Close()
	log.Println("Listening.")

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go s.handle(conn)
	}
}

func respond(w io.Writer, response string) {
	_, err := w.Write([]byte(response))
	if err != nil {
		log.Println(err)
	}

	// Scanner used in the client, which breaks on the newline
	w.Write([]byte("\n"))
}

func (s *Server) handle(rw io.ReadWriter) {
	reader := protocol.NewReader(rw)

	for {
		data, err := reader.Read()
		if err != nil {
			respond(rw, fmt.Sprint("Couldn't read message: ", err))
			return
		}

		msg, err := protocol.UnmarshalBinary(data, s.store.C)
		if err != nil {
			respond(rw, fmt.Sprint("Couldn't unmarshal binary: ", err))
			return
		}

		if msg.Cmd == protocol.Get {
			value, err := s.store.Get(msg.Key)
			if err != nil {
				respond(rw, fmt.Sprintf("GET: Error handling key '%s': %s", msg.Key, err.Error()))
				return
			}

			respond(rw, fmt.Sprintf("GET: %s", string(value)))
			_, err = rw.Write(value)
			if err != nil {
				log.Fatal(err)
			}
		} else if msg.Cmd == protocol.Set {
			expires, err := s.store.Set(msg.Key, string(msg.Data), msg.Expires)
			if err != nil {
				respond(rw, fmt.Sprintf("SET: Error setting key '%s': %s", msg.Key, err.Error()))
				return
			}

			respond(rw, fmt.Sprintf("Set '%s'. Expires: %s", msg.Key, expires))
		}
	}
}

type c struct{}

func (clock c) Now() time.Time {
	return time.Now()
}

func (clock c) Expired(t time.Time) bool {
	return t.Unix() < clock.Now().Unix()
}

func NewServer(cacheSize int) *Server {
	return &Server{store: cache.NewStore(uint64(cacheSize), c{})}
}
