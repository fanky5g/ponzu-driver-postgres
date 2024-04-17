package repository

import (
	"context"
	"fmt"
)

func (repo *repository) AutoMigrate(ctx context.Context) error {
	sqlString := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
	id UUID PRIMARY KEY,
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL,
	deleted_at TIMESTAMPTZ NULL,
	document jsonb NOT NULL
);
`, repo.model.Name())

	conn, err := repo.conn.Acquire(ctx)
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, sqlString)
	return err
}
