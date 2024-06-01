package repository

import (
	"context"
	"fmt"
	ponzuConstants "github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/entities"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

var DefaultQuerySize = 100

func (repo *repository) FindOneById(id string) (interface{}, error) {
	ctx := context.Background()
	conn, err := repo.conn.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()
	return repo.findOneByIdWithConn(id, conn)
}

func (repo *repository) Latest() (interface{}, error) {
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

func (repo *repository) FindOneBy(criteria map[string]interface{}) (interface{}, error) {
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

func (repo *repository) count(ctx context.Context, conn *pgxpool.Conn) (int, error) {
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

func (repo *repository) Find(order ponzuConstants.SortOrder, pagination *entities.Pagination) (int, []interface{}, error) {
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
	switch order {
	case ponzuConstants.Ascending:
		sortOrder = "ASC"
	}

	limit := DefaultQuerySize
	if pagination != nil && pagination.Count > 0 {
		limit = pagination.Count
	}

	sqlString := fmt.Sprintf(`
			SELECT id, created_at, updated_at, document
			FROM %s
			WHERE deleted_at IS NULL
			ORDER BY updated_at %s
			LIMIT %d
	`, repo.model.Name(), sortOrder, limit)

	if pagination != nil && pagination.Offset > 0 {
		sqlString = fmt.Sprintf(`
			%s
			OFFSET %d
`, sqlString, pagination.Offset)
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

func (repo *repository) FindAll() ([]interface{}, error) {
	var allResults []interface{}
	numResults, results, err := repo.Find(ponzuConstants.Descending, nil)
	if err != nil {
		return nil, err
	}

	if len(results) > 0 {
		allResults = append(allResults, results...)
	}

	fetched := len(results)
	for fetched != numResults {
		_, results, err = repo.Find(ponzuConstants.Descending, &entities.Pagination{
			Offset: fetched,
		})

		if len(results) > 0 {
			allResults = append(allResults, results...)
		}

		fetched = fetched + len(results)
	}

	return allResults, nil
}

func (repo *repository) findOneByIdWithConn(id string, conn *pgxpool.Conn) (interface{}, error) {
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
