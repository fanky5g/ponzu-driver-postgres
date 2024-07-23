package search

import (
	"errors"
	"fmt"

	"github.com/fanky5g/ponzu-driver-postgres/types"
)

type client struct {
	searchClients map[string]types.Search
}

// CreateIndex is a no-op as in postgres we don't have to create any indexes
func (c *client) CreateIndex(entityName string, entityType interface{}) error {
	return nil
}

func (c *client) GetIndex(entityName string) (types.Search, error) {
	if sc, ok := c.searchClients[entityName]; ok {
		return sc, nil
	}

	return nil, fmt.Errorf("%s search client not found", entityName)
}

func New(models []types.ModelInterface) (types.SearchClient, error) {
	searchClients := make(map[string]types.Search)
	for _, model := range models {
		entity := model.NewEntity()
		contentEntity, ok := entity.(types.Entity)
		if !ok {
			return nil, errors.New("entity must implement content.Entity interface")
		}

		var err error
		searchClients[contentEntity.EntityName()], err = NewEntitySearch(model)
		if err != nil {
			return nil, err
		}
	}

	return &client{searchClients: searchClients}, nil
}
