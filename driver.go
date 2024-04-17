package ponzu_driver_postgres

import (
	"context"
	"github.com/fanky5g/ponzu-driver-postgres/database"
	"github.com/fanky5g/ponzu-driver-postgres/repository"
	ponzuDriver "github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/models"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

type driver struct {
	conn         *pgxpool.Pool
	repositories map[string]ponzuDriver.Repository
}

func (database *driver) Get(token string) interface{} {
	if repo, ok := database.repositories[token]; ok {
		return repo
	}

	log.Panicf("Repository %s not found", token)
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

	repos := make(map[string]ponzuDriver.Repository)
	for _, model := range models {
		var repo ponzuDriver.Repository
		repo, err = repository.New(conn, model)
		if err != nil {
			return nil, err
		}

		repos[model.Name()] = repo
	}

	d := &driver{conn: conn, repositories: repos}

	return d, nil
}
