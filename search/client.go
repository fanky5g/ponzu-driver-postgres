package search

import "github.com/fanky5g/ponzu-driver-postgres/database"

type Client struct {
	db *database.Database
}

func New(db *database.Database) (*Client, error) {
	return &Client{db: db}, nil
}
