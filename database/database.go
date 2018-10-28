package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Conn represent a SQL connection.
type Conn struct {
	db *sql.DB
}

// New creates a new SQL connection.
func New(user, pass, name string) (*Conn, error) {
	connStr := fmt.Sprintf("%s:%s@/%s?parseTime=true", user, pass, name)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	return &Conn{db: db}, nil
}

// Close ensures that the connection is terminated.
func (c *Conn) Close() error {
	return c.db.Close()
}
