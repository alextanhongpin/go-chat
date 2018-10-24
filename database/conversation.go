package database

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
