package database

import "github.com/alextanhongpin/go-chat/entity"

func (c *Conn) AddFriend(userID, targetID, actorID int) error {
	_, err := c.db.Exec(`
	INSERT INTO friendship (user_id1, user_id2, actor_id, relationship) VALUES (?, ?, ?, (SELECT id FROM ref_relationship WHERE type = 'request'))`, userID, targetID, actorID)
	return err
}

func (c *Conn) RejectFriend(requestID int) error {
	_, err := c.db.Exec(`
	DELETE FROM friendship WHERE request_id = ?
	`, requestID)
	return err
}

func (c *Conn) AcceptFriend(requestID int) error {
	_, err := c.db.Exec(`
	UPDATE friendship SET relationship = (SELECT id FROM ref_relationship WHERE type = 'friend') WHERE id = ?
	`, requestID)
	return err

}
func (c *Conn) BlockFriend(requestID int) error {
	_, err := c.db.Exec(`
	UPDATE friendship SET relationship = (SELECT id FROM ref_relationship WHERE type = 'block') WHERE id = ?
	`, requestID)
	return err
}

func (c *Conn) GetRequestedFriends(id int) ([]entity.Friend, error) {
	rows, err := c.db.Query(`
		SELECT user_id2 AS id 
		FROM friendship 
		WHERE user_id1 = ? 
		AND actor_id = ? 
		AND relationship = (SELECT id FROM ref_relationship WHERE type = 'request') 
	`, id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.Friend
	for rows.Next() {
		var friend entity.Friend
		if err := rows.Scan(&friend.ID); err != nil {
			return nil, err
		}
		result = append(result, friend)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Conn) GetPendingFriends(id int) ([]entity.Friend, error) {
	rows, err := c.db.Query(`
		SELECT user_id2 AS id 
		FROM friendship 
		WHERE user_id1 = ? 
		AND actor_id <> ?
		AND relationship = (SELECT id FROM ref_relationship WHERE type = 'request') 
	`, id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.Friend
	for rows.Next() {
		var friend entity.Friend
		if err := rows.Scan(&friend.ID); err != nil {
			return nil, err
		}
		result = append(result, friend)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Conn) GetBlockedFriends(id int) ([]entity.Friend, error) {
	rows, err := c.db.Query(`
		SELECT user_id2 AS id 
		FROM friendship 
		WHERE user_id1 = ? 
		AND actor_id = ?
		AND relationship = (SELECT id FROM ref_relationship WHERE type = 'blocked') 
	`, id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.Friend
	for rows.Next() {
		var friend entity.Friend
		if err := rows.Scan(&friend.ID); err != nil {
			return nil, err
		}
		result = append(result, friend)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Conn) GetMutualFriends(id int) ([]entity.Friend, error) {
	rows, err := c.db.Query(`
		SELECT user_id2 AS id 
		FROM friendship 
		WHERE user_id1 = ? 
		AND relationship = (SELECT id FROM ref_relationship WHERE type = 'friend') 
	`, id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.Friend
	for rows.Next() {
		var friend entity.Friend
		if err := rows.Scan(&friend.ID); err != nil {
			return nil, err
		}
		result = append(result, friend)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
