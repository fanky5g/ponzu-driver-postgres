package ponzu_driver_postgres

import (
	"context"
	"fmt"
	"github.com/fanky5g/ponzu-driver-postgres/models"
	"github.com/fanky5g/ponzu-driver-postgres/repositories"
	ponzuContent "github.com/fanky5g/ponzu/content"
	ponzuDriver "github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/tokens"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"sync"
)

var (
	once sync.Once
	conn *pgxpool.Pool
)

type driver struct {
	conn *pgxpool.Pool
}

func (database *driver) Get(token tokens.Repository) interface{} {
	repository, err := repositories.New()
	if err != nil {
		log.Panicf("Failed to get repository %v: %v", token, err)
	}

	return repository
}

func (database *driver) Close() error {
	database.conn.Close()
	return nil
}

func New(
	contentTypes map[string]ponzuContent.Builder,
	contentModels map[string]*models.Model,
) (ponzuDriver.Database, error) {
	var err error
	once.Do(func() {
		var cfg *Config
		cfg, err = getConfig()
		if err != nil {
			err = fmt.Errorf("failed to get config: %v", err)
			return
		}

		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Name,
		)

		conn, err = pgxpool.New(context.Background(), dsn)
		if err != nil {
			err = fmt.Errorf("failed to connect to database: %v", err)
			return
		}
	})

	if err != nil {
		return nil, err
	}

	for _, model := range contentModels {
		models.RegisterModel(model)
	}

	if err = CreateTables(conn); err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	return &driver{conn: conn}, nil
}
