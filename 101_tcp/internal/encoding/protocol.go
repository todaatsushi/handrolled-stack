package encoding

// Message format:
// 8 bits

// Message format:
// 2 byte header
// version | lendata (2B)

const VERSION byte = 1
const HEADER_SIZE = 3
const MAX_PACKET_LEN = 10_000

type Basic struct{}

func (b Basic) Encode(msg string) []byte {
	return []byte("TODO")
}

func (b Basic) Decode(msg []byte) string {
	return "TODO"
}
