package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/fanky5g/ponzu-driver-postgres/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DefaultQuerySize = 100

func (repo *Repository) FindOneById(id string) (interface{}, error) {
	ctx := context.Background()
	conn, err := repo.conn.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()
	return repo.findOneByIdWithConn(id, conn)
}

func (repo *Repository) Latest() (interface{}, error) {
	sqlString := fmt.Sprintf(`
		SELECT id, created_at, updated_at, document
		FROM %s
		WHERE deleted_at IS NULL
		ORDER BY updated_at DESC
		LIMIT 1
`, repo.model.Name())

	ctx := context.Background()
	conn, err := repo.conn.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	return repo.ScanRow(conn.QueryRow(
		context.Background(),
		sqlString,
	))
}

func (repo *Repository) FindOneBy(criteria map[string]interface{}) (interface{}, error) {
	queryLength := len(criteria)
	whereClauses := make([]string, queryLength)
	values := make([]interface{}, queryLength)
	position := 0
	for field, value := range criteria {
		valueType, err := repo.valueType(value)
		if err != nil {
			return nil, err
		}

		whereClauses[position] = fmt.Sprintf(
			"(document->>'%s')::%s = $%d::%s",
			field,
			valueType,
			position+1,
			valueType,
		)

		values[position] = value
		position = position + 1
	}

	sqlString := fmt.Sprintf(
		"SELECT id, created_at, updated_at, document FROM %s WHERE %s AND deleted_at IS NULL",
		repo.model.Name(),
		strings.Join(whereClauses, " AND "),
	)

	ctx := context.Background()
	conn, err := repo.conn.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()
	return repo.ScanRow(conn.QueryRow(
		context.Background(),
		sqlString,
		values...,
	))
}

func (repo *Repository) count(ctx context.Context, conn *pgxpool.Conn) (int, error) {
	sqlString := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s
		WHERE deleted_at IS NULL
`, repo.model.Name())

	count := 0
	err := conn.QueryRow(ctx, sqlString).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *Repository) Find(order string, pagination types.Paginator) (int, []interface{}, error) {
	ctx := context.Background()
	conn, err := repo.conn.Acquire(ctx)
	if err != nil {
		return 0, nil, err
	}

	defer conn.Release()
	count, err := repo.count(ctx, conn)
	if err != nil {
		return 0, nil, err
	}

	if count == 0 {
		return 0, nil, nil
	}

	sortOrder := "DESC"
	switch strings.ToLower(order) {
	case "asc":
		sortOrder = "ASC"
	}

	limit := DefaultQuerySize
	if pagination != nil && pagination.GetCount() > 0 {
		limit = pagination.GetCount()
	}

	sqlString := fmt.Sprintf(`
			SELECT id, created_at, updated_at, document
			FROM %s
			WHERE deleted_at IS NULL
			ORDER BY updated_at %s
			LIMIT %d
	`, repo.model.Name(), sortOrder, limit)

	if pagination != nil && pagination.GetOffSet() > 0 {
		sqlString = fmt.Sprintf(`
			%s
			OFFSET %d
`, sqlString, pagination.GetOffSet())
	}

	rows, err := conn.Query(ctx, sqlString)
	if err != nil {
		return 0, nil, err
	}

	defer rows.Close()
	results := make([]interface{}, 0)
	for rows.Next() {
		var result interface{}
		if result, err = repo.ScanRow(rows); err != nil {
			return 0, nil, err
		}

		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return 0, nil, err
	}

	return count, results, nil
}

func (repo *Repository) FindAll() ([]interface{}, error) {
	var allResults []interface{}
	numResults, results, err := repo.Find("DESC", nil)
	if err != nil {
		return nil, err
	}

	if len(results) > 0 {
		allResults = append(allResults, results...)
	}

	fetched := len(results)
	for fetched != numResults {
		_, results, err = repo.Find("DESC", newPaginator(0, fetched))

		if len(results) > 0 {
			allResults = append(allResults, results...)
		}

		fetched = fetched + len(results)
	}

	return allResults, nil
}

func (repo *Repository) findOneByIdWithConn(id string, conn *pgxpool.Conn) (interface{}, error) {
	sqlString := fmt.Sprintf(
		"SELECT id, created_at, updated_at, document FROM %s WHERE id = $1::uuid AND deleted_at IS NULL",
		repo.model.Name(),
	)

	return repo.ScanRow(conn.QueryRow(
		context.Background(),
		sqlString,
		id,
	))
}

type p struct {
	count  int
	offset int
}

func (p p) GetOffSet() int {
	return p.offset
}

func (p p) GetCount() int {
	return p.count
}

func newPaginator(count, offset int) types.Paginator {
	return p{count, offset}
}
