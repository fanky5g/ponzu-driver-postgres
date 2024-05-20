package repository

import (
	"context"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	conn  *pgxpool.Pool
	model models.ModelInterface
}

type Repository interface {
	driver.Repository
	RowScanner
}

func New(conn *pgxpool.Pool, model models.ModelInterface) (Repository, error) {
	repo := &repository{
		conn:  conn,
		model: model,
	}

	if err := repo.AutoMigrate(context.Background()); err != nil {
		return nil, err
	}

	return repo, nil
}
