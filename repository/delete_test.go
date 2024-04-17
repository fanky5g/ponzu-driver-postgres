package repository

import (
	"context"
	"fmt"
	"github.com/fanky5g/ponzu-driver-postgres/database"
	ponzuConstants "github.com/fanky5g/ponzu/constants"
	ponzuDriver "github.com/fanky5g/ponzu/driver"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DeleteTestSuite struct {
	suite.Suite
	repo ponzuDriver.Repository
	conn *pgxpool.Pool
}

func (s *DeleteTestSuite) SetupSuite() {
	conn, err := database.GetConnection(context.Background())
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

func (s *DeleteTestSuite) TearDownTest() {
	ctx := context.Background()
	_, err := s.conn.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", testModelToken))
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *DeleteTestSuite) TestDeleteById() {
	entity := &testEntity{
		Name:  "Foo Bar",
		Email: "foo@bar.domain",
		Age:   39,
	}

	ins, err := s.repo.Insert(entity)
	if err != nil {
		s.T().Fatal(err)
	}

	insertId := ins.(*testEntity).ID
	item, err := s.repo.FindOneById(insertId)
	if err != nil {
		s.T().Fatal(err)
	}

	if item == nil {
		s.T().Fatalf("Expected item to be inserted correctly but was nil")
	}

	err = s.repo.DeleteById(insertId)
	if err != nil {
		s.T().Fatalf("Expected err to be nil. Got: %v", err)
	}

	item, err = s.repo.FindOneById(insertId)
	if assert.NoError(s.T(), err) {
		assert.Nil(s.T(), item)
	}
}

func (s *DeleteTestSuite) TestDeleteByFieldEqual() {
	entity := &testEntity{
		Name:  "Foo Bar",
		Email: "foo@bar.domain",
		Age:   39,
	}

	ins, err := s.repo.Insert(entity)
	if err != nil {
		s.T().Fatal(err)
	}

	insertId := ins.(*testEntity).ID
	item, err := s.repo.FindOneById(insertId)
	if err != nil {
		s.T().Fatal(err)
	}

	if item == nil {
		s.T().Fatalf("Expected item to be inserted correctly but was nil")
	}

	err = s.repo.DeleteBy("email", ponzuConstants.Equal, entity.Email)
	if err != nil {
		s.T().Fatalf("Expected err to be nil. Got: %v", err)
	}

	item, err = s.repo.FindOneById(insertId)
	if assert.NoError(s.T(), err) {
		assert.Nil(s.T(), item)
	}
}

func (s *DeleteTestSuite) TestDeleteByFieldLessThan() {
	for _, entity := range testEntities {
		_, err := s.repo.Insert(entity)
		if err != nil {
			s.T().Fatal(err)
		}
	}

	err := s.repo.DeleteBy("age", ponzuConstants.LessThan, 20)
	if err != nil {
		s.T().Fatalf("Expected err to be nil. Got: %v", err)
	}

	matches, err := s.repo.FindAll()
	if assert.NoError(s.T(), err) {
		assert.Len(s.T(), matches, 2)
	}
}

func (s *DeleteTestSuite) TestDeleteByFieldGreaterThan() {
	for _, entity := range testEntities {
		_, err := s.repo.Insert(entity)
		if err != nil {
			s.T().Fatal(err)
		}
	}

	err := s.repo.DeleteBy("age", ponzuConstants.GreaterThan, 10)
	if err != nil {
		s.T().Fatalf("Expected err to be nil. Got: %v", err)
	}

	matches, err := s.repo.FindAll()
	if assert.NoError(s.T(), err) {
		assert.Len(s.T(), matches, 1)
	}
}

func (s *DeleteTestSuite) TestDeleteByFieldGreaterThanOrEqualTo() {
	for _, entity := range testEntities {
		_, err := s.repo.Insert(entity)
		if err != nil {
			s.T().Fatal(err)
		}
	}

	err := s.repo.DeleteBy("age", ponzuConstants.GreaterThanOrEqualTo, 10)
	if err != nil {
		s.T().Fatalf("Expected err to be nil. Got: %v", err)
	}

	matches, err := s.repo.FindAll()
	if assert.NoError(s.T(), err) {
		assert.Len(s.T(), matches, 0)
	}
}

func (s *DeleteTestSuite) TestDeleteByFieldLessThanOrEqualTo() {
	for _, entity := range testEntities {
		_, err := s.repo.Insert(entity)
		if err != nil {
			s.T().Fatal(err)
		}
	}

	err := s.repo.DeleteBy("age", ponzuConstants.LessThanOrEqualTo, 20)
	if err != nil {
		s.T().Fatalf("Expected err to be nil. Got: %v", err)
	}

	matches, err := s.repo.FindAll()
	if assert.NoError(s.T(), err) {
		assert.Len(s.T(), matches, 1)
	}
}

func TestDelete(t *testing.T) {
	suite.Run(t, new(DeleteTestSuite))
}
