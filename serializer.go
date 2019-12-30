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
	Serialize(ch chan *EntityQueue, w io.Writer) error
	Deserialize(r io.Reader) chan *EntityQueue
}

// JSONSerializer is simple line json serializer
type JSONSerializer struct{}

// NewJSONSerializer is constructor of JSONSerializer
func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{}
}

// Serialize of JSONSerializer marshals BadEntity to JSON and append line feed at tail.
func (x *JSONSerializer) Serialize(ch chan *EntityQueue, w io.Writer) error {
	for q := range ch {
		if q.Error != nil {
			return q.Error
		}

		for _, e := range q.Entities {
			raw, err := json.Marshal(e)
			if err != nil {
				return errors.Wrapf(err, "Fail to marshal entity: %v", e)
			}

			line := append(raw, []byte("\n")...)
			if _, err := w.Write(line); err != nil {
				return errors.Wrapf(err, "Fail to write entity: %v", e)
			}
		}
	}

	return nil
}

// Deserialize of JSONSerializer reads reader and unmarshal nd-json.
func (x *JSONSerializer) Deserialize(r io.Reader) chan *EntityQueue {
	ch := make(chan *EntityQueue, jsonSerializerBufSize)
	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			var entity BadEntity
			raw := scanner.Bytes()
			if err := json.Unmarshal(raw, &entity); err != nil {
				ch <- &EntityQueue{
					Error: errors.Wrapf(err, "Fail to unmarshal serialized entity as json: %s", string(raw)),
				}
				return
			}

			ch <- &EntityQueue{Entities: []*BadEntity{&entity}}
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
func (x *GzipJSONSerializer) Serialize(ch chan *EntityQueue, w io.Writer) error {
	writer := gzip.NewWriter(w)

	for q := range ch {
		if q.Error != nil {
			return q.Error
		}

		for _, e := range q.Entities {
			raw, err := json.Marshal(e)
			if err != nil {
				return errors.Wrapf(err, "Fail to marshal entity: %v", e)
			}

			line := append(raw, []byte("\n")...)
			if _, err := writer.Write(line); err != nil {
				return errors.Wrapf(err, "Fail to write entity: %v", e)
			}
		}
	}

	if err := writer.Close(); err != nil {
		return errors.Wrap(err, "Fail to close GzipJSONSerializer writer")
	}

	return nil
}

// Deserialize of GzipJSONSerializer reads reader and unmarshal gzipped nd-json.
func (x *GzipJSONSerializer) Deserialize(r io.Reader) chan *EntityQueue {
	ch := make(chan *EntityQueue, jsonSerializerBufSize)
	go func() {
		defer close(ch)
		reader, err := gzip.NewReader(r)
		if err != nil {
			ch <- &EntityQueue{
				Error: errors.Wrapf(err, "Fail to create a new reader for GzipJSONSerializer"),
			}
			return
		}

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			var entity BadEntity
			raw := scanner.Bytes()
			if err := json.Unmarshal(raw, &entity); err != nil {
				ch <- &EntityQueue{
					Error: errors.Wrapf(err, "Fail to unmarshal serialized entity as json: %s", string(raw)),
				}
				return
			}

			ch <- &EntityQueue{Entities: []*BadEntity{&entity}}
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
func (x *MsgpackSerializer) Serialize(ch chan *EntityQueue, w io.Writer) error {
	enc := msgpack.NewEncoder(w)

	for q := range ch {
		if q.Error != nil {
			return q.Error
		}

		for _, e := range q.Entities {
			if err := enc.Encode(e); err != nil {
				return errors.Wrapf(err, "Fail to encode entity by MsgpackSerializer: %v", e)
			}
		}
	}

	return nil
}

// Deserialize of MsgpackSerializer reads reader and unmarshal gzipped nd-json.
func (x *MsgpackSerializer) Deserialize(r io.Reader) chan *EntityQueue {
	ch := make(chan *EntityQueue, jsonSerializerBufSize)

	go func() {
		defer close(ch)
		dec := msgpack.NewDecoder(r)

		for {
			var entity BadEntity
			err := dec.Decode(&entity)
			if err == io.EOF {
				return
			} else if err != nil {
				ch <- &EntityQueue{
					Error: errors.Wrapf(err, "Fail to decode msgpack format"),
				}
				return
			}

			ch <- &EntityQueue{Entities: []*BadEntity{&entity}}
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
func (x *GzipMsgpackSerializer) Serialize(ch chan *EntityQueue, w io.Writer) error {
	writer := gzip.NewWriter(w)
	enc := msgpack.NewEncoder(writer)

	for q := range ch {
		if q.Error != nil {
			return q.Error
		}

		for _, e := range q.Entities {
			if err := enc.Encode(e); err != nil {
				return errors.Wrapf(err, "Fail to encode entity by GzipMsgpackSerializer: %v", e)
			}
		}
	}

	if err := writer.Close(); err != nil {
		return errors.Wrap(err, "Fail to close GzipMsgpackSerializer writer")
	}

	return nil
}

// Deserialize of GzipMsgpackSerializer reads reader and unmarshal gzipped nd-json.
func (x *GzipMsgpackSerializer) Deserialize(r io.Reader) chan *EntityQueue {
	ch := make(chan *EntityQueue, jsonSerializerBufSize)

	go func() {
		defer close(ch)
		reader, err := gzip.NewReader(r)
		if err != nil {
			ch <- &EntityQueue{
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
				ch <- &EntityQueue{
					Error: errors.Wrapf(err, "Fail to decode msgpack format"),
				}
				return
			}

			ch <- &EntityQueue{Entities: []*BadEntity{&entity}}
		}
	}()

	return ch
}
