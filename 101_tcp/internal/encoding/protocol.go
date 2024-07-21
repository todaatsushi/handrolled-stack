package encoding

import (
	"encoding/binary"
	"errors"
)

// Message format:
// 2 byte header
// version | lendata (2B)

const VERSION byte = 1
const HEADER_SIZE = 3
const MAX_PACKET_LEN = 10_000

type Basic struct{}

func (b Basic) Encode(msg string) []byte {
	data := []byte(msg)

	// Build header
	lenData := uint16(len(data))
	lenDataHeader := make([]byte, 2)
	binary.BigEndian.PutUint16(lenDataHeader, lenData)

	// Assemble
	packet := []byte{}
	packet = append(packet, VERSION)
	packet = append(packet, lenDataHeader...)
	packet = append(packet, data...)
	return packet
}

func (b Basic) Decode(msg []byte) (string, error) {
	if len(msg) < HEADER_SIZE {
		return "", errors.New("Not enough data.")
	}

	header := msg[:HEADER_SIZE]
	version := header[0]
	if version != VERSION {
		return "", errors.New("Version mismatch.")
	}

	lenDataBytes := header[1:HEADER_SIZE]
	lenData := binary.BigEndian.Uint16(lenDataBytes)

	if lenData != 0 && len(msg) <= HEADER_SIZE {
		return "", errors.New("Data length specified but no data attached.")
	}

	data := msg[HEADER_SIZE:]
	if len(data) != int(lenData) {
		return "", errors.New("Data size doesn't match length specified.")
	}

	return string(data), nil
}
