package badman

import "io"

// Repository test interfaces
func RepositoryPut(repo Repository, entity BadEntity) error {
	return repo.put(entity)
}
func RepositoryGet(repo Repository, name string) ([]BadEntity, error) {
	return repo.get(name)
}
func RepositoryDel(repo Repository, name string) error {
	return repo.del(name)
}
func RepositoryDump(repo Repository) chan *BadEntityMessage {
	return repo.dump()
}
func RepositoryClear(repo Repository) error {
	return repo.clear()
}

func SerializerSerialize(sre Serializer, ch chan *BadEntityMessage, w io.Writer) error {
	return sre.serialize(ch, w)
}

func SerializerDeserialize(sre Serializer, r io.Reader) chan *BadEntityMessage {
	return sre.deserialize(r)
}
