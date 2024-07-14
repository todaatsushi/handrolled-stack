package encoding

type Translator interface {
	Encode() []byte
	Decode() string
}
