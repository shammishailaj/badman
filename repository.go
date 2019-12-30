package badman

import (
	"time"
)

// BadEntity is IP address or domain name that is appeared in BlackList. Name indicates both IP address and domain name.
type BadEntity struct {
	Name    string
	SavedAt time.Time
	Src     string
	Reason  string // optional
}

// Repository is interface of data store.
type Repository interface {
	Put(entities []*BadEntity) error
	Get(name string) ([]BadEntity, error)
	Del(name string) error
	Dump() chan *EntityQueue
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

func (x *inMemoryRepository) Put(entities []*BadEntity) error {
	for i := range entities {
		srcMap, ok := x.data[entities[i].Name]
		if !ok {
			srcMap = make(map[string]BadEntity)
			x.data[entities[i].Name] = srcMap
		}
		srcMap[entities[i].Src] = *entities[i]
	}
	return nil
}

func (x *inMemoryRepository) Get(name string) ([]BadEntity, error) {
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

func (x *inMemoryRepository) Del(name string) error {
	delete(x.data, name)
	return nil
}

func (x *inMemoryRepository) Dump() chan *EntityQueue {
	ch := make(chan *EntityQueue)
	go func() {
		defer close(ch)
		for _, srcMap := range x.data {
			var q EntityQueue
			for _, entity := range srcMap {
				q.Entities = append(q.Entities, &entity)
			}
			ch <- &q
		}
	}()
	return ch
}
