package encoding

import (
	"encoding/binary"
	"errors"
	"net"
)

type MessageReader struct {
	conn    net.Conn
	buffer  []byte
	scratch []byte
	init    bool
}

func NewMessageReader(conn net.Conn) *MessageReader {
	return &MessageReader{
		conn:    conn,
		buffer:  []byte{},
		scratch: make([]byte, 1024),
		init:    true,
	}
}

func (r *MessageReader) getDataSize(data []byte) (int, error) {
	if len(data) < HEADER_SIZE {
		return -1, errors.New("Not enough data.")
	}

	// Last 2 bytes of the header is the length of the data
	dataLength := int(binary.BigEndian.Uint16(data[1:HEADER_SIZE]))
	return dataLength, nil
}

func (r *MessageReader) hasCompletePacket(data []byte) bool {
	dataLength, err := r.getDataSize(data)
	if err != nil {
		return false
	}
	return len(data) >= HEADER_SIZE+dataLength
}

func (r *MessageReader) HasData() bool {
	return len(r.buffer) > 0 || r.init
}

func (r *MessageReader) Read() (data []byte, err error) {
	for {
		if len(r.buffer) > MAX_PACKET_LEN {
			return []byte{}, errors.New("Packet exceeds max length.")
		}

		if r.hasCompletePacket(r.buffer) {
			size, err := r.getDataSize(r.buffer)
			if err != nil {
				return []byte{}, err
			}
			parsed := r.buffer[:size+HEADER_SIZE]
			r.buffer = r.buffer[size+HEADER_SIZE:]
			r.init = false
			return parsed, nil
		}

		n, err := r.conn.Read(r.scratch)
		if err != nil {
			return []byte{}, err
		}

		r.buffer = append(r.buffer, r.scratch[:n]...)
	}
}
