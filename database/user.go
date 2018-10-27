package database

type User struct {
	ID   string
	Name string
}

type UserRepository interface {
	GetUser(id string) (User, error)
	GetUserByName(name string) (User, error)
}

func (c *Conn) GetUser(id string) (user User, err error) {
	err = c.db.QueryRow("SELECT * FROM user WHERE id = ?", id).Scan(
		&user.ID,
		&user.Name,
	)
	return
}

func (c *Conn) GetUserByName(name string) (user User, err error) {
	err = c.db.QueryRow("SELECT * FROM user WHERE name = ?", name).Scan(
		&user.ID,
		&user.Name,
	)
	return
}
