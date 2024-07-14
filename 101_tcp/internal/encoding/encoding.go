package encoding

type Translator interface {
	Encode(msg string) []byte
	Decode(msg []byte) string
}
