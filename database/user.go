package database

import "github.com/alextanhongpin/go-chat/entity"

func (c *Conn) GetUser(id string) (user entity.User, err error) {
	err = c.db.QueryRow("SELECT * FROM user WHERE id = ?", id).Scan(
		&user.ID,
		&user.Name,
	)
	return
}

func (c *Conn) GetUserByName(name string) (user entity.User, err error) {
	err = c.db.QueryRow("SELECT * FROM user WHERE name = ?", name).Scan(
		&user.ID,
		&user.Name,
	)
	return
}
