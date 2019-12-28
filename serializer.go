package badman

import "io"

// Serializer converts array of BadEntity to byte array and the reverse.
type Serializer interface {
	Serialize(ch chan *BadEntityMessage, w io.Writer) error
	Deserialize(r io.Reader) chan *BadEntityMessage
}
