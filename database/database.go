package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Conn struct {
	db *sql.DB
}

func New(user, pass, name string) (*Conn, error) {
	connStr := fmt.Sprintf("%s:%s@/%s", user, pass, name)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	return &Conn{db: db}, nil
}

func (c *Conn) Close() error {
	return c.db.Close()
}
