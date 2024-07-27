package repository

import (
	"context"
	"github.com/fanky5g/ponzu/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	conn  *pgxpool.Pool
	model models.ModelInterface
}

func (repo *Repository) TableName() string {
    return repo.model.Name() 
}

func New(conn *pgxpool.Pool, model models.ModelInterface) (*Repository, error) {
	repo := &Repository{
		conn:  conn,
		model: model,
	}

	if err := repo.AutoMigrate(context.Background()); err != nil {
		return nil, err
	}

	return repo, nil
}
