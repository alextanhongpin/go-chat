package database

import (
	"database/sql"
	"strconv"

	"github.com/alextanhongpin/go-chat/entity"
)

// GetUser returns a user by id.
func (c *Conn) GetUser(id string) (user entity.User, err error) {
	err = c.db.QueryRow("SELECT id, name FROM user WHERE id = ?", id).Scan(
		&user.ID,
		&user.Name,
	)
	return
}

// GetUserByName returns a user from the given name or an error if the user
// doesn't exist.
func (c *Conn) GetUserByName(name string) (user entity.User, err error) {
	err = c.db.QueryRow("SELECT id, name FROM user WHERE name = ?", name).Scan(
		&user.ID,
		&user.Name,
	)
	return
}

func (c *Conn) GetUserByEmail(email string) (user entity.User, err error) {
	err = c.db.QueryRow(`
		SELECT 
			id, 
			name, 
			email, 
			created_at, 
			updated_at, 
			deleted_at, 
			hashed_password 
		FROM 
			user 
		WHERE 
			email = ?`,
		email).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt, &user.HashedPassword)
	if err == sql.ErrNoRows {
		err = entity.ErrUserNotFound
	}
	return
}

func (c *Conn) CreateUser(user *entity.User) error {
	result, err := c.db.Exec("INSERT INTO user (name, email, hashed_password) VALUES (?, ?, ?)",
		user.Name, user.Email, user.HashedPassword)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = strconv.FormatInt(id, 10)
	return nil
}
