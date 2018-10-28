package database

import "github.com/alextanhongpin/go-chat/entity"

func (c *Conn) CreateConversationReply(userID, roomID string, text string) (int64, error) {
	stmt, err := c.db.Prepare("INSERT INTO conversation (user_id, room_id, text) VALUES (?, ?, ?)")
	if err != nil {
		return -1, err
	}
	res, err := stmt.Exec(userID, roomID, text)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

func (c *Conn) GetConversations(roomID string) ([]entity.Conversation, error) {
	stmt, err := c.db.Prepare("SELECT user_id, created_at, text FROM conversation WHERE room_id = ? ORDER BY created_at DESC LIMIT 10")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []entity.Conversation
	for rows.Next() {
		var res entity.Conversation
		err := rows.Scan(&res.UserID, &res.CreatedAt, &res.Text)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	return result, nil
}
