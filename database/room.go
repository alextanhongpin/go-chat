package database

import (
	"log"

	"github.com/alextanhongpin/go-chat/entity"
)

// CreateRoom create a new room.
func (c *Conn) CreateRoom(users ...string) error {
	tx, err := c.db.Begin()

	// Create a room where the conversation will happen.
	res, err := tx.Exec("INSERT INTO room VALUES (NULL)")
	if err != nil {
		tx.Rollback()
		return err
	}
	roomID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// Create a many to many relationship between user and room, a junction table.
	for _, userID := range users {
		_, err := tx.Exec("INSERT INTO user_room (user_id, room_id) VALUES (?, ?)", roomID, userID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// GetRoom returns the roomID for the given userID.
func (c *Conn) GetRoom(userID string) ([]string, error) {
	rows, err := c.db.Query("SELECT (room_id) FROM user_room WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []string
	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}

// GetRooms returns a list of user in the room.
func (c *Conn) GetRooms(userID string) ([]entity.UserRoom, error) {
	log.Println("getting rooms", userID)
	stmt := `SELECT user_id, room_id, name FROM (SELECT user_id, room_id FROM user_room WHERE room_id IN (SELECT room_id FROM user_room WHERE user_id = ?) AND user_id <> ?) a INNER JOIN user b ON b.id = a.user_id`

	rows, err := c.db.Query(stmt, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.UserRoom
	for rows.Next() {
		var userRoom entity.UserRoom
		err := rows.Scan(&userRoom.UserID, &userRoom.RoomID, &userRoom.Name)
		if err != nil {
			return nil, err
		}
		result = append(result, userRoom)
	}
	return result, nil
}
