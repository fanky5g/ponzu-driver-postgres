package models

import (
	"github.com/fanky5g/ponzu-driver-postgres/utils"
	"time"
)

var Models map[string]interface{}

// This package encompasses ponzu system models. Content type models live in the application space and are generated
// alongside ponzu types.

type Model struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func RegisterModel(mm ...interface{}) {
	if Models == nil {
		Models = make(map[string]interface{})
	}

	for _, m := range mm {
		Models[utils.TableName(m)] = m
	}
}
