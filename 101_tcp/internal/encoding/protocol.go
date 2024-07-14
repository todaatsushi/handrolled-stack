package encoding

import "errors"

// Message format:
// 8 bits

// vv lendata
// 00 000000 data

const VERSION = 1

type Basic struct {
	decoded string
	encoded []byte
}

func New(encoded []byte, decoded string) (Basic, error) {
	if len(encoded) == 0 && len(decoded) == 0 {
		return Basic{}, errors.New("Must provide either encoded or decoded.")
	}
	if len(encoded) > 0 && len(decoded) > 0 {
		return Basic{}, errors.New("Must provide one of encoded or decoded.")
	}

	return Basic{
		decoded: decoded,
		encoded: encoded,
	}, nil
}

func (b Basic) Encode() []byte {
	return []byte("TODO")
}

func (b Basic) Decode() string {
	return "TODO"
}
