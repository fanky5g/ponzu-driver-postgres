package search

import (
	"context"
	"fmt"
	"github.com/fanky5g/ponzu-driver-postgres/connection"
	"github.com/fanky5g/ponzu-driver-postgres/database/repository"
	ponzuDriver "github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type searchClient struct {
	searchableFields []string
	tableName        string

	conn *pgxpool.Pool
	repo repository.RowScanner
}

func (s *searchClient) Search(query string, limit, offset int) ([]interface{}, error) {
	results, _, err := s.SearchWithPagination(query, limit, offset)
	return results, err
}

func (s *searchClient) SearchWithPagination(query string, limit, offset int) ([]interface{}, int, error) {
	ctx := context.Background()
	conn, err := s.conn.Acquire(ctx)
	if err != nil {
		return nil, 0, err
	}

	defer conn.Release()

	queryLength := len(s.searchableFields)
	whereClauses := make([]string, queryLength)
	values := make([]interface{}, queryLength)
	position := 0
	for _, field := range s.searchableFields {
		whereClauses[position] = fmt.Sprintf(
			"(document->>'%s') LIKE $1",
			field,
		)

		values[position] = query
		position = position + 1
	}

	size := repository.DefaultQuerySize
	if limit > 0 {
		size = limit
	}

	whereClause := strings.Join(whereClauses, " OR ")

	sqlString := fmt.Sprintf(`
			SELECT id, created_at, updated_at, document
			FROM %s
			WHERE (%s) AND deleted_at IS NULL
			ORDER BY updated_at DESC
			LIMIT %d
	`, s.tableName, whereClause, size)

	if offset > 0 {
		sqlString = fmt.Sprintf(`
			%s
			OFFSET %d
`, sqlString, offset)
	}

	value := fmt.Sprintf("%%%s%%", query)
	count, err := s.count(ctx, conn, whereClause, value)
	if err != nil {
		return nil, 0, err
	}

	rows, err := conn.Query(ctx, sqlString, value)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()
	results := make([]interface{}, 0)
	for rows.Next() {
		var result interface{}
		if result, err = s.repo.ScanRow(rows); err != nil {
			return nil, 0, err
		}

		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return results, count, nil
}

func (s *searchClient) count(ctx context.Context, conn *pgxpool.Conn, whereClause, value string) (int, error) {
	sqlString := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s
		WHERE (%s) AND deleted_at IS NULL
`, s.tableName, whereClause)

	count := 0
	err := conn.QueryRow(ctx, sqlString, value).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Update is a no-op as with postgres we don't have to actually update any index
func (s *searchClient) Update(id string, data interface{}) error {
	return nil
}

// Delete is a no-op as with postgres we don't have to delete from an index.
//
//	Deletion from database from repository should be sufficient
func (s *searchClient) Delete(id string) error {
	return nil
}

func NewEntitySearch(model models.ModelInterface) (ponzuDriver.SearchInterface, error) {
	conn, err := connection.Get(context.Background())
	if err != nil {
		return nil, err
	}

	return NewEntitySearchWithConn(conn, model)
}

func NewEntitySearchWithConn(conn *pgxpool.Pool, model models.ModelInterface) (ponzuDriver.SearchInterface, error) {
	repo, err := repository.New(conn, model)
	if err != nil {
		return nil, err
	}

	searchableFields, err := getSearchableFields(model.NewEntity())
	if err != nil {
		return nil, err
	}

	return &searchClient{
		searchableFields: searchableFields,
		conn:             conn,
		repo:             repo,
		tableName:        model.Name(),
	}, nil
}
