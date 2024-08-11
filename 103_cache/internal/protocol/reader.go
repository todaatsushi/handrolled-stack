package protocol

import (
	"encoding/binary"
	"errors"
	"io"
)

type DataReader struct {
	stream  io.Reader
	buf     []byte
	scratch []byte
}

func NewReader(r io.Reader) *DataReader {
	return &DataReader{
		stream:  r,
		buf:     []byte{},
		scratch: make([]byte, 1024),
	}
}

func (r *DataReader) lenData() (int, error) {
	if len(r.buf) < HEADER_SIZE {
		return -1, errors.New("Not enough data read.")
	}

	dataLenBytes := r.buf[6:8]
	return int(binary.BigEndian.Uint16(dataLenBytes)), nil
}

func (r *DataReader) lenKey() (int, error) {
	if len(r.buf) < HEADER_SIZE {
		return -1, errors.New("Not enough data read.")
	}

	keyLenBytes := r.buf[4:6]
	return int(binary.BigEndian.Uint16(keyLenBytes)), nil
}

func (r *DataReader) lenTotal() (int, error) {
	lenKey, err := r.lenKey()
	if err != nil {
		return -1, err
	}

	lenData, err := r.lenData()
	if err != nil {
		return -1, err
	}

	return HEADER_SIZE + lenData + lenKey, nil
}

func (r *DataReader) complete() bool {
	lenTotal, err := r.lenTotal()
	if err != nil {
		return false
	}
	return lenTotal == len(r.buf)
}

func (r *DataReader) Read() (data []byte, err error) {
	ptr := 0
	for {
		isComplete := r.complete()
		if isComplete {
			lenTotal, err := r.lenTotal()
			if err != nil {
				return []byte{}, err
			}

			message := r.buf[:lenTotal]
			r.buf = r.buf[lenTotal:]
			return message, nil
		}

		numRead, err := r.stream.Read(r.scratch)
		if err != nil {
			return []byte{}, err
		}

		toCopy := r.scratch[ptr : ptr+numRead]
		r.buf = append(r.buf, toCopy...)

		ptr += numRead
	}
}
