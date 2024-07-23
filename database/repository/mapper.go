package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fanky5g/ponzu-driver-postgres/types"
	"github.com/google/uuid"
)

// MapToEntity maps a domain model to a Database Model
func (repo *Repository) MapToEntity(entity interface{}) *types.Model {
	model := &types.Model{
		Document: repo.model.ToDocument(entity),
	}

	if identifiable, ok := entity.(types.Identifiable); ok {
		if identifiable.ItemID() == "" {
			identifiable.SetItemID(uuid.New().String())
		}

		model.ID = uuid.Must(uuid.Parse(identifiable.ItemID()))
	} else {
		model.ID = uuid.New()
	}

	if temporal, ok := entity.(types.Temporal); ok {
		if temporal.CreatedAt() != 0 {
			model.CreatedAt = time.Unix(temporal.CreatedAt(), 0)
		}

		if temporal.UpdatedAt() != 0 {
			model.UpdatedAt = time.Unix(temporal.UpdatedAt(), 0)
		}
	}

	return model
}

func (repo *Repository) MapFromEntity(model *types.Model) (interface{}, error) {
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

	if identifiable, ok := entity.(types.Identifiable); ok {
		identifiable.SetItemID(model.ID.String())
	}

	if temporal, ok := entity.(types.Temporal); ok {
		temporal.SetCreatedAt(model.CreatedAt)
		temporal.SetUpdatedAt(model.UpdatedAt)
	}

	return entity, nil
}
