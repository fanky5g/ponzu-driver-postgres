package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
    "github.com/fanky5g/ponzu-driver-postgres/types"
)

type Repository struct {
	conn  *pgxpool.Pool
	model types.ModelInterface
}

func New(conn *pgxpool.Pool, model types.ModelInterface) (*Repository, error) {
	repo := &Repository{
		conn:  conn,
		model: model,
	}

	if err := repo.AutoMigrate(context.Background()); err != nil {
		return nil, err
	}

	return repo, nil
}
