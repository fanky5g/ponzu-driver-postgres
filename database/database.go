package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/fanky5g/ponzu-driver-postgres/connection"
	"github.com/fanky5g/ponzu-driver-postgres/database/repository"
	"github.com/fanky5g/ponzu-driver-postgres/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRepositoryNotFound = errors.New("repository not found")

type Driver struct {
	conn         *pgxpool.Pool
	repositories map[string]*repository.Repository
}

func New(models ...types.ModelInterface) (*Driver, error) {
	ctx := context.Background()
	conn, err := connection.Get(ctx)

	if err != nil {
		return nil, err
	}

	repos := make(map[string]*repository.Repository)
	for _, model := range models {
		entity := model.NewEntity()
		persistable, ok := entity.(types.Persistable)
		if !ok {
			return nil, fmt.Errorf("entity %T is not persistable", entity)
		}

		var repo *repository.Repository
		repo, err = repository.New(conn, model)
		if err != nil {
			return nil, err
		}

		repos[persistable.GetRepositoryToken()] = repo
	}

	d := &Driver{conn: conn, repositories: repos}

	return d, nil
}
