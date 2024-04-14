package ponzu_driver_postgres

import (
	"context"
	"fmt"
	"github.com/fanky5g/ponzu-driver-postgres/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func createTable(ctx context.Context, conn *pgxpool.Pool, tableName string) error {
	sqlString := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
	id UUID PRIMARY KEY,
	created_at DATE NOT NULL,
	updated_at DATE NOT NULL,
	document jsonb NOT NULL
);
`, tableName)

	_, err := conn.Exec(ctx, sqlString)
	return err
}

func CreateTables(conn *pgxpool.Pool) error {
	ctx := context.Background()

	for tableName := range models.Models {
		if err := createTable(ctx, conn, tableName); err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
	}

	return nil
}
