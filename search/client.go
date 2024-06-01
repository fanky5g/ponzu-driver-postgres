package search

import (
	"errors"
	"fmt"
	"github.com/fanky5g/ponzu/content"
	ponzuDriver "github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/models"
)

type client struct {
	searchClients map[string]ponzuDriver.SearchInterface
}

// CreateIndex is a no-op as in postgres we don't have to create any indexes
func (c *client) CreateIndex(entityName string, entityType interface{}) error {
	return nil
}

func (c *client) GetIndex(entityName string) (ponzuDriver.SearchInterface, error) {
	if sc, ok := c.searchClients[entityName]; ok {
		return sc, nil
	}

	return nil, fmt.Errorf("%s search client not found", entityName)
}

func New(models []models.ModelInterface) (ponzuDriver.SearchClientInterface, error) {
	searchClients := make(map[string]ponzuDriver.SearchInterface)
	for _, model := range models {
		entity := model.NewEntity()
		contentEntity, ok := entity.(content.Entity)
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
