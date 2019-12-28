package badman

import (
	"time"
)

// BadEntity is IP address or domain name that is appeared in BlackList. Name indicates both IP address and domain name.
type BadEntity struct {
	Name    string
	SavedAt time.Time
	Src     string
}

// BadEntityMessage is messaging format for Dump() to send both of error and
type BadEntityMessage struct {
	Error  error
	Entity *BadEntity
}

// Repository is interface of data store.
type Repository interface {
	put(entity BadEntity) error
	get(name string) ([]BadEntity, error)
	del(name string) error
	dump() chan *BadEntityMessage
	clear() error
}

// inMemoryRepository is in-memory type repository.
type inMemoryRepository struct {
	data map[string]map[string]BadEntity
}

// NewInMemoryRepository is constructor of inMemoryRepository
func NewInMemoryRepository() Repository {
	repo := &inMemoryRepository{}
	repo.init()
	return repo
}

func (x *inMemoryRepository) init() {
	x.data = make(map[string]map[string]BadEntity)
}

func (x *inMemoryRepository) put(entity BadEntity) error {
	srcMap, ok := x.data[entity.Name]
	if !ok {
		srcMap = make(map[string]BadEntity)
		x.data[entity.Name] = srcMap
	}
	srcMap[entity.Src] = entity
	return nil
}

func (x *inMemoryRepository) get(name string) ([]BadEntity, error) {
	srcMap, ok := x.data[name]
	if !ok {
		return nil, nil
	}

	var entities []BadEntity
	for _, entity := range srcMap {
		entities = append(entities, entity)
	}

	return entities, nil
}

func (x *inMemoryRepository) del(name string) error {
	delete(x.data, name)
	return nil
}

func (x *inMemoryRepository) dump() chan *BadEntityMessage {
	ch := make(chan *BadEntityMessage)
	go func() {
		defer close(ch)
		for _, srcMap := range x.data {
			for _, entity := range srcMap {
				ch <- &BadEntityMessage{Entity: &entity}
			}
		}
	}()
	return ch
}

func (x *inMemoryRepository) clear() error {
	x.init()
	return nil
}
