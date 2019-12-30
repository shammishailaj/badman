package badman

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

const jsonSerializerBufSize = 32

// Serializer converts array of BadEntity to byte array and the reverse.
type Serializer interface {
	Serialize(ch chan *BadEntityMessage, w io.Writer) error
	Deserialize(r io.Reader) chan *BadEntityMessage
}

// JSONSerializer is simple line json serializer
type JSONSerializer struct{}

// NewJSONSerializer is constructor of JSONSerializer
func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{}
}

// Serialize of JSONSerializer marshals BadEntity to JSON and append line feed at tail.
func (x *JSONSerializer) Serialize(ch chan *BadEntityMessage, w io.Writer) error {
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

// Deserialize of JSONSerializer reads reader and unmarshal nd-json.
func (x *JSONSerializer) Deserialize(r io.Reader) chan *BadEntityMessage {
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

// GzipJSONSerializer is simple line json serializer
type GzipJSONSerializer struct{}

// NewGzipJSONSerializer is constructor of GzipJSONSerializer
func NewGzipJSONSerializer() *GzipJSONSerializer {
	return &GzipJSONSerializer{}
}

// Serialize of GzipJSONSerializer marshals BadEntity to gzipped JSON and append line feed at tail.
func (x *GzipJSONSerializer) Serialize(ch chan *BadEntityMessage, w io.Writer) error {
	writer := gzip.NewWriter(w)

	for msg := range ch {
		if msg.Error != nil {
			return msg.Error
		}

		raw, err := json.Marshal(msg.Entity)
		if err != nil {
			return errors.Wrapf(err, "Fail to marshal entity: %v", msg.Entity)
		}

		line := append(raw, []byte("\n")...)
		if _, err := writer.Write(line); err != nil {
			return errors.Wrapf(err, "Fail to write entity: %v", msg.Entity)
		}
	}

	if err := writer.Close(); err != nil {
		return errors.Wrap(err, "Fail to close GzipJSONSerializer writer")
	}

	return nil
}

// Deserialize of GzipJSONSerializer reads reader and unmarshal gzipped nd-json.
func (x *GzipJSONSerializer) Deserialize(r io.Reader) chan *BadEntityMessage {
	ch := make(chan *BadEntityMessage, jsonSerializerBufSize)
	go func() {
		defer close(ch)
		reader, err := gzip.NewReader(r)
		if err != nil {
			ch <- &BadEntityMessage{
				Error: errors.Wrapf(err, "Fail to create a new reader for GzipJSONSerializer"),
			}
			return
		}

		scanner := bufio.NewScanner(reader)
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
