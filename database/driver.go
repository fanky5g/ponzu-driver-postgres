package database

import (
	"github.com/fanky5g/ponzu-driver-postgres/database/repository"
	"github.com/fanky5g/ponzu-driver-postgres/types"
)

func (database *Driver) Insert(repositoryToken string, entity interface{}) (interface{}, error) {
	repo, err := database.getRepository(repositoryToken)
	if err != nil {
		return nil, err
	}

	return repo.Insert(entity)
}

func (database *Driver) Latest(repositoryToken string) (interface{}, error) {
	repo, err := database.getRepository(repositoryToken)
	if err != nil {
		return nil, err
	}

	return repo.Latest()
}

func (database *Driver) UpdateById(repositoryToken string, id string, update interface{}) (interface{}, error) {
	repo, err := database.getRepository(repositoryToken)
	if err != nil {
		return nil, err
	}

	return repo.UpdateById(id, update)
}

func (database *Driver) Find(repositoryToken, order string, pagination types.Paginator) (int, []interface{}, error) {
	repo, err := database.getRepository(repositoryToken)
	if err != nil {
		return 0, nil, err
	}

	return repo.Find(order, pagination)
}

func (database *Driver) FindOneById(repositoryToken, id string) (interface{}, error) {
	repo, err := database.getRepository(repositoryToken)
	if err != nil {
		return nil, err
	}

	return repo.FindOneById(id)
}

func (database *Driver) FindOneBy(repositoryToken string, criteria map[string]interface{}) (interface{}, error) {
	repo, err := database.getRepository(repositoryToken)
	if err != nil {
		return nil, err
	}

	return repo.FindOneBy(criteria)
}

func (database *Driver) FindAll(repositoryToken string) ([]interface{}, error) {
	repo, err := database.getRepository(repositoryToken)
	if err != nil {
		return nil, err
	}

	return repo.FindAll()
}

func (database *Driver) DeleteById(repositoryToken, id string) error {
	repo, err := database.getRepository(repositoryToken)
	if err != nil {
		return err
	}

	return repo.DeleteById(id)
}

func (database *Driver) DeleteBy(repositoryToken, field string, operator string, value interface{}) error {
	repo, err := database.getRepository(repositoryToken)
	if err != nil {
		return err
	}

	return repo.DeleteBy(field, operator, value)
}

func (database *Driver) Close() error {
	database.conn.Close()
	return nil
}

func (database *Driver) getRepository(repositoryToken string) (*repository.Repository, error) {
	repo, ok := database.repositories[repositoryToken]
	if !ok {
		return nil, ErrRepositoryNotFound
	}

	return repo, nil
}
