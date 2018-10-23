package database

import "github.com/alextanhongpin/go-chat/entity"

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

func (c *Conn) GetRoom(userID string) ([]int64, error) {
	rows, err := c.db.Query("SELECT (room_id) FROM user_room WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []int64
	for rows.Next() {
		var i int64
		err := rows.Scan(&i)
		if err != nil {
			return nil, err
		}
		result = append(result, i)
	}
	return result, nil
}

func (c *Conn) GetRooms(userID string) ([]entity.UserRoom, error) {
	stmt := `SELECT l.user_id, l.room_id FROM user_room l LEFT JOIN user_room r ON l.room_id = r.room_id WHERE r.user_id = ? AND r.user_id != l.user_id`
	rows, err := c.db.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.UserRoom
	for rows.Next() {
		var userRoom entity.UserRoom
		err := rows.Scan(&userRoom.UserID, &userRoom.RoomID)
		if err != nil {
			return nil, err
		}
		result = append(result, userRoom)
	}
	return result, nil
}
