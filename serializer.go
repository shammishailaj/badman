package badman

import (
	"bufio"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

const jsonSerializerBufSize = 32

// Serializer converts array of BadEntity to byte array and the reverse.
type Serializer interface {
	serialize(ch chan *BadEntityMessage, w io.Writer) error
	deserialize(r io.Reader) chan *BadEntityMessage
}

// JSONSerializer is simple line json serializer
type JSONSerializer struct{}

// NewJSONSerializer is constructor of JSONSerializer
func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{}
}

func (x *JSONSerializer) serialize(ch chan *BadEntityMessage, w io.Writer) error {
	for msg := range ch {
		if msg.Error != nil {
			return msg.Error
		}

		raw, err := json.Marshal(msg.Entity)
		if err != nil {
			return errors.Wrapf(err, "Fail to marshal entity: %v", msg.Entity)
		}

		line := append(raw, []byte("\n")...)
		if _, err := w.Write(line); err != nil {
			return errors.Wrapf(err, "Fail to write entity: %v", msg.Entity)
		}
	}

	return nil
}

func (x *JSONSerializer) deserialize(r io.Reader) chan *BadEntityMessage {
	ch := make(chan *BadEntityMessage, jsonSerializerBufSize)
	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			var entity BadEntity
			raw := scanner.Bytes()
			if err := json.Unmarshal(raw, &entity); err != nil {
				ch <- &BadEntityMessage{
					Error: errors.Wrapf(err, "Fail to unmarshal serialized entity as json: %s", string(raw)),
				}
				return
			}

			ch <- &BadEntityMessage{Entity: &entity}
		}
	}()

	return ch
}
