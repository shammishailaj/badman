package badman

import (
	"time"
)

// BadEntity is IP address or domain name that is appeared in BlackList. Name indicates both IP address and domain name.
type BadEntity struct {
	Name     string
	SavedAt  time.Time
	Src      string
	IsActive bool
}

// BadEntityMessage is messaging format for Dump() to send both of error and
type BadEntityMessage struct {
	Error  error
	Entity *BadEntity
}

// Repository is interface of data store.
type Repository interface {
	put(entity BadEntity) error
	get(name string) (*BadEntity, error)
	dump() chan *BadEntityMessage
}

// inMemoryRepository is in-memory type repository.
type inMemoryRepository struct {
	data map[string]*BadEntity
}

// NewInMemoryRepository is constructor of inMemoryRepository
func NewInMemoryRepository() Repository {
	repo := inMemoryRepository{
		data: make(map[string]*BadEntity),
	}
	return &repo
}

func (x *inMemoryRepository) put(entity BadEntity) error {
	return nil
}

func (x *inMemoryRepository) get(name string) (*BadEntity, error) {
	return nil, nil
}

func (x *inMemoryRepository) dump() chan *BadEntityMessage {
	return nil
}
