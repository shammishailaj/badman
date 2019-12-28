package badman

import "io"

// BadMan is Main interface of badman pacakge.
type BadMan struct {
	repo Repository
	ser  Serializer
}

// New is constructor of BadMan
func New() *BadMan {
	return &BadMan{}
}

// Insert adds an entity one by one. It's expected to use adding IoC by feed or something like that.
func (x *BadMan) Insert(entity BadEntity) error {
	return nil
}

// Lookup searches BadEntity (both of IP address and domain name). If not found, the function returns (nil, nil).
func (x *BadMan) Lookup(name string) (*BadEntity, error) {
	return nil, nil
}

// Download accesses blacklist data via Sources and store entities that is included in blacklist into repository.
func (x *BadMan) Download(srcSet []Source) error {
	return nil
}

// Dump output serialized data into w to save current repository.
func (x *BadMan) Dump(w io.Writer) error {
	if err := x.ser.Serialize(x.repo.dump(), w); err != nil {
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

		if err := x.repo.put(*msg.Entity); err != nil {
			return err
		}
	}
	return nil
}

// -----------------------------------
// Utilities

// ReplaceRepository changes Repository to store entities. Entities in old repository are copied to new repository before replacing.
func (x *BadMan) ReplaceRepository(repo Repository) error {
	for msg := range x.repo.dump() {
		if msg.Error != nil {
			return msg.Error
		}

		if err := repo.put(*msg.Entity); err != nil {
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
