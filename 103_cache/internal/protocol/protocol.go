package protocol

import "time"

// Version (1B) | Command (1B) | TTL (1B) | Length (2B) | Data (x)
const VERSION byte = 1

type Command byte

const (
	_ Command = iota
	Get
	Set
	Update
)

type Message struct {
	Cmd     Command
	Data    []byte
	Expires time.Time
}
