package database

import "github.com/alextanhongpin/go-chat/entity"

// GetUser returns a user by id.
func (c *Conn) GetUser(id string) (user entity.User, err error) {
	err = c.db.QueryRow("SELECT * FROM user WHERE id = ?", id).Scan(
		&user.ID,
		&user.Name,
	)
	return
}

// GetUserByName returns a user from the given name or an error if the user
// doesn't exist.
func (c *Conn) GetUserByName(name string) (user entity.User, err error) {
	err = c.db.QueryRow("SELECT * FROM user WHERE name = ?", name).Scan(
		&user.ID,
		&user.Name,
	)
	return
}
