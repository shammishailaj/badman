package badman

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
	msgpack "github.com/vmihailenco/msgpack/v4"
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

// MsgpackSerializer is MessagePack serializer
type MsgpackSerializer struct{}

// NewMsgpackSerializer is constructor of MsgpackSerializer
func NewMsgpackSerializer() *MsgpackSerializer {
	return &MsgpackSerializer{}
}

// Serialize of MsgpackSerializer encodes BadEntity to MessagePack format.
func (x *MsgpackSerializer) Serialize(ch chan *BadEntityMessage, w io.Writer) error {
	enc := msgpack.NewEncoder(w)

	for msg := range ch {
		if msg.Error != nil {
			return msg.Error
		}

		if err := enc.Encode(msg.Entity); err != nil {
			return errors.Wrapf(err, "Fail to encode entity by MsgpackSerializer: %v", msg.Entity)
		}
	}

	return nil
}

// Deserialize of MsgpackSerializer reads reader and unmarshal gzipped nd-json.
func (x *MsgpackSerializer) Deserialize(r io.Reader) chan *BadEntityMessage {
	ch := make(chan *BadEntityMessage, jsonSerializerBufSize)

	go func() {
		defer close(ch)
		dec := msgpack.NewDecoder(r)

		for {
			var entity BadEntity
			err := dec.Decode(&entity)
			if err == io.EOF {
				return
			} else if err != nil {
				ch <- &BadEntityMessage{
					Error: errors.Wrapf(err, "Fail to decode msgpack format"),
				}
				return
			}

			ch <- &BadEntityMessage{Entity: &entity}
		}
	}()

	return ch
}

// GzipMsgpackSerializer is MessagePack serializer
type GzipMsgpackSerializer struct{}

// NewGzipMsgpackSerializer is constructor of GzipMsgpackSerializer
func NewGzipMsgpackSerializer() *GzipMsgpackSerializer {
	return &GzipMsgpackSerializer{}
}

// Serialize of GzipMsgpackSerializer encodes BadEntity to MessagePack format.
func (x *GzipMsgpackSerializer) Serialize(ch chan *BadEntityMessage, w io.Writer) error {
	writer := gzip.NewWriter(w)
	enc := msgpack.NewEncoder(writer)

	for msg := range ch {
		if msg.Error != nil {
			return msg.Error
		}

		if err := enc.Encode(msg.Entity); err != nil {
			return errors.Wrapf(err, "Fail to encode entity by GzipMsgpackSerializer: %v", msg.Entity)
		}
	}

	if err := writer.Close(); err != nil {
		return errors.Wrap(err, "Fail to close GzipMsgpackSerializer writer")
	}

	return nil
}

// Deserialize of GzipMsgpackSerializer reads reader and unmarshal gzipped nd-json.
func (x *GzipMsgpackSerializer) Deserialize(r io.Reader) chan *BadEntityMessage {
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

		dec := msgpack.NewDecoder(reader)

		for {
			var entity BadEntity
			err := dec.Decode(&entity)
			if err == io.EOF {
				return
			} else if err != nil {
				ch <- &BadEntityMessage{
					Error: errors.Wrapf(err, "Fail to decode msgpack format"),
				}
				return
			}

			ch <- &BadEntityMessage{Entity: &entity}
		}
	}()

	return ch
}
