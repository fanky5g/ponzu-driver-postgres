package repository

import (
	"context"
	"fmt"
	"github.com/fanky5g/ponzu-driver-postgres/connection"
	ponzuDriver "github.com/fanky5g/ponzu/driver"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type InsertTestSuite struct {
	suite.Suite
	repo ponzuDriver.Repository
	conn *pgxpool.Pool
}

func (s *InsertTestSuite) SetupSuite() {
	conn, err := connection.Get(context.Background())
	if err != nil {
		s.T().Fatal(err)
	}

	repo, err := New(conn, new(testModel))
	if err != nil {
		s.FailNow(err.Error())
	}

	s.repo = repo
	s.conn = conn
}

func (s *InsertTestSuite) TearDownSuite() {
	ctx := context.Background()
	_, err := s.conn.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", testModelToken))
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *InsertTestSuite) TestInsert() {
	entity := &testEntity{
		Name:  "Foo Bar",
		Email: "foo@bar.domain",
		Age:   39,
	}

	ins, err := s.repo.Insert(entity)
	if assert.NoError(s.T(), err) && assert.NotNil(s.T(), ins) {
		assert.IsType(s.T(), ins, new(testEntity))
		saved := ins.(*testEntity)

		assert.Equal(s.T(), saved.Name, entity.Name)
		assert.Equal(s.T(), saved.Email, entity.Email)
		assert.Equal(s.T(), saved.Age, entity.Age)
	}
}

func TestInsert(t *testing.T) {
	suite.Run(t, new(InsertTestSuite))
}
