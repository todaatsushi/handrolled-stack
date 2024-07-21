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

func (b Basic) Decode(msg []byte) string {
	return "TODO"
}
