package ponzu_driver_postgres

import (
	"context"
	"fmt"
	"github.com/fanky5g/ponzu-driver-postgres/database"
	"github.com/fanky5g/ponzu-driver-postgres/repository"
	ponzuDriver "github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/models"
	"github.com/fanky5g/ponzu/tokens"
	"github.com/jackc/pgx/v5/pgxpool"
)

type driver struct {
	conn         *pgxpool.Pool
	repositories map[tokens.RepositoryToken]ponzuDriver.Repository
}

func (database *driver) GetRepositoryByToken(token tokens.RepositoryToken) ponzuDriver.Repository {
	if repo, ok := database.repositories[token]; ok {
		return repo
	}

	return nil
}

func (database *driver) Close() error {
	database.conn.Close()
	return nil
}

func New(models []models.ModelInterface) (ponzuDriver.Database, error) {
	ctx := context.Background()
	conn, err := database.GetConnection(ctx)

	if err != nil {
		return nil, err
	}

	repos := make(map[tokens.RepositoryToken]ponzuDriver.Repository)
	for _, model := range models {
		entity := model.NewEntity()
		persistable, ok := entity.(entities.Persistable)
		if !ok {
			return nil, fmt.Errorf("entity %T is not persistable", entity)
		}

		var repo ponzuDriver.Repository
		repo, err = repository.New(conn, model)
		if err != nil {
			return nil, err
		}

		repos[persistable.GetRepositoryToken()] = repo
	}

	d := &driver{conn: conn, repositories: repos}

	return d, nil
}
