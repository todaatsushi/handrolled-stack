package server

import "github.com/todaatsushi/handrolled-cache/internal/server"

func Start(port, cacheSize int) error {
	s := server.NewServer(cacheSize)
	return s.Run(port)
}
