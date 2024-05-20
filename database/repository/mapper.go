package repository

import (
	"encoding/json"
	"fmt"
	ponzuItem "github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/models"
	"github.com/google/uuid"
	"time"
)

// MapToEntity maps a domain model to a Database Model
func (repo *repository) MapToEntity(entity interface{}) *models.Model {
	model := &models.Model{
		Document: repo.model.ToDocument(entity),
	}

	if identifiable, ok := entity.(ponzuItem.Identifiable); ok {
		if identifiable.ItemID() == "" {
			identifiable.SetItemID(uuid.New().String())
		}

		model.ID = uuid.Must(uuid.Parse(identifiable.ItemID()))
	} else {
		model.ID = uuid.New()
	}

	if temporal, ok := entity.(ponzuItem.Temporal); ok {
		if temporal.CreatedAt() != 0 {
			model.CreatedAt = time.Unix(temporal.CreatedAt(), 0)
		}

		if temporal.UpdatedAt() != 0 {
			model.UpdatedAt = time.Unix(temporal.UpdatedAt(), 0)
		}
	}

	return model
}

func (repo *repository) MapFromEntity(model *models.Model) (interface{}, error) {
	entity := repo.model.NewEntity()

	if model.Document != nil {
		doc, err := json.Marshal(model.Document)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(doc, entity); err != nil {
			return nil, fmt.Errorf("failed to unmarshal document to %T: %v", entity, err)
		}
	}

	if identifiable, ok := entity.(ponzuItem.Identifiable); ok {
		identifiable.SetItemID(model.ID.String())
	}

	if temporal, ok := entity.(ponzuItem.Temporal); ok {
		temporal.SetCreatedAt(model.CreatedAt)
		temporal.SetUpdatedAt(model.UpdatedAt)
	}

	return entity, nil
}
