package badman

import (
	"io"

	"github.com/pkg/errors"
)

// BadMan is Main interface of badman pacakge.
type BadMan struct {
	repo Repository
	ser  Serializer
}

// New is constructor of BadMan
func New() *BadMan {
	return &BadMan{
		repo: NewInMemoryRepository(),
		ser:  NewGzipMsgpackSerializer(),
	}
}

// Insert adds an entity one by one. It's expected to use adding IoC by feed or something like that.
func (x *BadMan) Insert(entity BadEntity) error {
	return x.repo.Put([]*BadEntity{&entity})
}

// Lookup searches BadEntity (both of IP address and domain name). If not found, the function returns ([]BadEntity{}, nil). A reason to return list of BadEntity is that multiple blacklists may have same entity.
func (x *BadMan) Lookup(name string) ([]BadEntity, error) {
	return x.repo.Get(name)
}

// Download accesses blacklist data via Sources and store entities that is included in blacklist into repository.
func (x *BadMan) Download(srcSet []Source) error {
	msgCh := make(chan *EntityQueue, 128)

	for i := 0; i < len(srcSet); i++ {
		src := srcSet[i]

		go func() {
			for msg := range src.Download() {
				msgCh <- msg
			}

			// Send empty message to notify termination
			defer func() { msgCh <- &EntityQueue{} }()
		}()
	}

	closed := 0
	for q := range msgCh {
		if len(q.Entities) == 0 && q.Error == nil {
			closed++
			if closed >= len(srcSet) {
				break
			}
			continue
		}

		if q.Error != nil {
			return errors.Wrap(q.Error, "Fail to download from source")
		}
		if err := x.repo.Put(q.Entities); err != nil {
			return errors.Wrapf(err, "Fail to put downloaded entity: %v", q.Entities)
		}
	}

	return nil
}

// Dump output serialized data into w to save current repository.
func (x *BadMan) Dump(w io.Writer) error {
	if err := x.ser.Serialize(x.repo.Dump(), w); err != nil {
		return err
	}

	return nil
}

// Load input data that is serialized by Dump(). Please note to use same Serializer for Dump and Load.
func (x *BadMan) Load(r io.Reader) error {
	for msg := range x.ser.Deserialize(r) {
		if msg.Error != nil {
			return msg.Error
		}

		if err := x.repo.Put(msg.Entities); err != nil {
			return err
		}
	}
	return nil
}

// -----------------------------------
// Utilities

// ReplaceRepository changes Repository to store entities. Entities in old repository are copied to new repository before replacing.
func (x *BadMan) ReplaceRepository(repo Repository) error {
	for msg := range x.repo.Dump() {
		if msg.Error != nil {
			return msg.Error
		}

		if err := repo.Put(msg.Entities); err != nil {
			return err
		}
	}

	x.repo = repo
	return nil
}

// ReplaceSerializer just changes Serializer with ser.
func (x BadMan) ReplaceSerializer(ser Serializer) {
	x.ser = ser
}
