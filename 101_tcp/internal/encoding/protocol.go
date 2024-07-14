package encoding

// Message format:
// 8 bits

// vv lendata
// 00 000000 data

const VERSION = 1

type Basic struct{}

func (b Basic) Encode(msg string) []byte {
	return []byte("TODO")
}

func (b Basic) Decode(msg []byte) string {
	return "TODO"
}
