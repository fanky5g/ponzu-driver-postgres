package types

import (
	"github.com/google/uuid"
	"time"
    "reflect"
)

type Paginator interface {
	GetCount() int
	GetOffSet() int
}

type DocumentInterface interface {
	Value() (interface{}, error)
	Scan(src interface{}) error
}

type ModelInterface interface {
	Name() string
	ToDocument(entity interface{}) DocumentInterface
	NewEntity() interface{}
}

type Search interface {
	Update(id string, data interface{}) error
	Delete(id string) error
	Search(query string, count, offset int) ([]interface{}, error)
	SearchWithPagination(query string, count, offset int) ([]interface{}, int, error)
}

type SearchClient interface {
	CreateIndex(entityName string, entityType interface{}) error
	GetIndex(entityName string) (Search, error)
}

type CustomizableSearchAttributes interface {
	GetSearchableAttributes() map[string]reflect.Type
}

type Model struct {
	ID        uuid.UUID         `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeletedAt time.Time         `json:"deleted_at"`
	Document  DocumentInterface `json:"document"`
}

type Entity interface {
	EntityName() string
}

type Persistable interface {
	GetRepositoryToken() string
}

type Identifiable interface {
	ItemID() string
	SetItemID(string)
}

type Temporal interface {
	CreatedAt() int64
	SetCreatedAt(time.Time)
	UpdatedAt() int64
	SetUpdatedAt(time.Time)
}
