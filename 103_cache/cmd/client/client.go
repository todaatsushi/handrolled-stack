package client

import "github.com/todaatsushi/handrolled-cache/internal/client"

func Start(port int) error {
	return client.Dial(port)
}
