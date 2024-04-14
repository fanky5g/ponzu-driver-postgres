package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Document struct{}

func (document *Document) Value() (driver.Value, error) {
	return json.Marshal(document)
}

func (document *Document) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed")
	}

	err := json.Unmarshal(source, &document)
	if err != nil {
		return err
	}

	return nil
}
