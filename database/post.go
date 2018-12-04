package database

import (
	"time"

	"github.com/alextanhongpin/go-chat/entity"
)

func (c *Conn) GetPosts() ([]entity.Post, error) {
	rows, err := c.db.Query(`
		SELECT id, text, created_at 
		FROM post
		WHERE deleted_at = '1900-01-01 00:00:00'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.Post
	for rows.Next() {
		var post entity.Post
		err := rows.Scan(
			&post.ID,
			&post.Text,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, post)
	}
	return result, nil
}

func (c *Conn) GetPost(id string) (entity.Post, error) {
	var post entity.Post
	err := c.db.QueryRow(`
		SELECT id, text, created_at 
		FROM post 
		WHERE id = ? 
		AND deleted_at <> '1900-01-01 00:00:00'
	`, id).Scan(
		&post.ID,
		&post.Text,
		&post.CreatedAt,
	)
	return post, err
}

func (c *Conn) CreatePost(post entity.Post) (int64, error) {
	result, err := c.db.Exec(`
		INSERT INTO post (text, user_id) 
		VALUES (?, ?)
	`, post.Text, post.UserID)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (c *Conn) UpdatePost(post entity.Post) error {
	_, err := c.db.Exec(`
		UPDATE post SET text = ? 
		WHERE id = ?
	`, post.Text, post.ID)
	return err
}

func (c *Conn) DeletePost(id string) error {
	_, err := c.db.Exec(`
		UPDATE post SET deleted_at = ?
		WHERE id = ?
	`, time.Now(), id)
	return err
}
