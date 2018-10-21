package database

type User struct {
	ID   int
	Name string
}

type UserRepository interface {
	GetUser(id string) (User, error)
	GetUserByName(name string) (User, error)
}

func (c *Conn) GetUser(id string) (User, error) {
	var user User
	err := c.db.QueryRow("SELECT * FROM user WHERE id = ?", id).Scan(
		&user.ID,
		&user.Name,
	)
	return user, err
}

func (c *Conn) GetUserByName(name string) (User, error) {
	var user User
	err := c.db.QueryRow("SELECT * FROM user WHERE name = ?", name).Scan(
		&user.ID,
		&user.Name,
	)
	return user, err
}
