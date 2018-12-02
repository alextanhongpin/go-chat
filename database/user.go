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

func (c *Conn) GetUsers(id int) ([]entity.Friend, error) {
	rows, err := c.db.Query(`
		SELECT 
			u.id, 
			u.name, 
			COALESCE(ref.type, ""), 
			COALESCE(fr.actor_id, "") 
		FROM user u 
		LEFT JOIN (
			SELECT * FROM friendship 
			WHERE user_id1 = ? 
			OR user_id2 = ?) fr 
			ON fr.user_id1 = u.id 
			OR fr.user_id2 = u.id 
		LEFT JOIN ref_relationship ref 
		ON ref.id = fr.relationship 
		WHERE u.id <> ?;
	`, id, id, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.Friend
	for rows.Next() {
		var actorID string
		var user entity.Friend
		if err := rows.Scan(&user.ID, &user.Name, &user.Status, &actorID); err != nil {
			return nil, err
		}
		user.IsRequested = strconv.Itoa(id) == actorID
		result = append(result, user)
	}
	return result, rows.Err()
}
