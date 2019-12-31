package badman

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/pkg/errors"
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

type dynamoRepository struct {
	table dynamo.Table
	msgCh chan *EntityQueue
}

type dynamoEntityItem struct {
	Name    string    `dynamo:"name"`
	Src     string    `dynamo:"src"`
	SavedAt time.Time `dynamo:"saved_at"`
	Reason  string    `dynamo:"reason"`
}

const dynamoBatchSize = 25

// NewDynamoRepository is constructor of dynamoRepository
func NewDynamoRepository(region, tableName string) Repository {
	ssn := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	db := dynamo.New(ssn)

	repo := &dynamoRepository{
		table: db.Table(tableName),
		msgCh: make(chan *EntityQueue),
	}

	return repo
}

func (x *dynamoRepository) flush(items []interface{}) error {
	wrote, err := x.table.Batch().Write().Put(items...).Run()
	if err != nil {
		return errors.Wrapf(err, "Fail to flush DynamoRepository: %v", items)
	}
	if wrote != len(items) {
		return errors.Wrapf(err, "Invalid wrote item number, expect %d but actual %d", len(items), wrote)
	}

	return nil
}

func (x *dynamoRepository) Put(entities []*BadEntity) error {
	var items []interface{}

	for i := 0; i < len(entities); i++ {

		item := dynamoEntityItem{
			Name:    entities[i].Name,
			Src:     entities[i].Src,
			SavedAt: entities[i].SavedAt,
			Reason:  entities[i].Reason,
		}
		items = append(items, item)

		if len(items) >= dynamoBatchSize {
			if err := x.flush(items); err != nil {
				return err
			}
			items = []interface{}{}
		}

	}

	if len(items) > 0 {
		if err := x.flush(items); err != nil {
			return err
		}
	}

	return nil
}

func (x *dynamoRepository) Get(name string) ([]BadEntity, error) {
	var items []dynamoEntityItem
	err := x.table.Get("name", name).All(&items)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to get entities from DynamoDB: %s", name)
	}

	var entities []BadEntity
	for _, item := range items {
		entities = append(entities, BadEntity{
			Name:    item.Name,
			Src:     item.Src,
			SavedAt: item.SavedAt,
			Reason:  item.Reason,
		})
	}

	return entities, nil
}

func (x *dynamoRepository) Del(name string) error {
	var items []dynamoEntityItem
	if err := x.table.Get("name", name).All(&items); err != nil {
		return errors.Wrapf(err, "Fail to get entities for deleteItems from DynamoDB: %s", name)
	}

	var keys []dynamo.Keyed
	for _, item := range items {
		keys = append(keys, &dynamo.Keys{item.Name, item.Src})
	}

	if wrote, err := x.table.Batch("name", "src").Write().Delete(keys...).Run(); err != nil {
		return errors.Wrapf(err, "Fail to delete entity from DynamoDB: %s (%v)", name, keys)
	} else if wrote != len(keys) {
		return errors.Wrapf(err, "Invalid delete item number, expect %d but actual %d", len(keys), wrote)
	}

	return nil
}

func (x *dynamoRepository) Dump() chan *EntityQueue {
	// Dump is not supported for DynamoDB because Scan requries massive resource against blacklist.
	return nil
}
