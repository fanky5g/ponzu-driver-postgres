package repositories

import (
	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/entities"
)

type repository struct {
}

func (r *repository) Insert(entity interface{}) (interface{}, error) {
	return nil, nil
}

func (r *repository) Latest() (interface{}, error) {
	return nil, nil
}

func (r *repository) UpdateById(id string, update interface{}) (interface{}, error) {
	return nil, nil
}

func (r *repository) Search(search *entities.Search) (int, []interface{}, error) {
	return 0, nil, nil
}

func (r *repository) FindOneById(id string) (interface{}, error) {
	return nil, nil
}

func (r *repository) FindOneBy(criteria map[string]interface{}) (interface{}, error) {
	return nil, nil
}

func (r *repository) FindAll() ([]interface{}, error) {
	return nil, nil
}

func (r *repository) DeleteById(id string) error {
	return nil
}

func (r *repository) DeleteBy(field string, operator constants.ComparisonOperator, value interface{}) error {
	return nil
}

func New() (driver.Repository, error) {
	return &repository{}, nil
}
