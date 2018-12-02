package database

import "github.com/alextanhongpin/go-chat/entity"

func (c *Conn) AddFriend(userID, targetID, actorID int) error {
	_, err := c.db.Exec(`
	INSERT INTO friendship (user_id1, user_id2, actor_id, relationship) VALUES (?, ?, ?, (SELECT id FROM ref_relationship WHERE type = 'request'))`, userID, targetID, actorID)
	return err
}

func (c *Conn) RejectFriend(a, b int) error {
	_, err := c.db.Exec(`
	DELETE FROM friendship WHERE user_id1 = ? AND user_id2 = ?
	`, a, b)
	return err
}

func (c *Conn) AcceptFriend(a, b int) error {
	_, err := c.db.Exec(`
	UPDATE friendship 
	SET relationship = (SELECT id FROM ref_relationship WHERE type = 'friend') 
	WHERE user_id1 = ? AND user_id2 = ?
	`, a, b)
	return err

}
func (c *Conn) BlockFriend(a, b int) error {
	_, err := c.db.Exec(`
	UPDATE friendship 
	SET relationship = (SELECT id FROM ref_relationship WHERE type = 'block') 
	WHERE user_id1 = ? AND user_id2 = ? 
	`, a, b)
	return err
}

func (c *Conn) GetRequestedFriends(id int) ([]entity.Friend, error) {
	rows, err := c.db.Query(`
		SELECT user_id2 AS id 
		FROM friendship 
		AND actor_id = ? 
		AND relationship = (SELECT id FROM ref_relationship WHERE type = 'request') 
	`, id)
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
	return result, rows.Err()
}

func (c *Conn) GetPendingFriends(id int) ([]entity.Friend, error) {
	rows, err := c.db.Query(`
		SELECT user_id2 AS id 
		FROM friendship 
		WHERE user_id1 = ? OR user_id2 = ? 
		AND actor_id <> ?
		AND relationship = (SELECT id FROM ref_relationship WHERE type = 'request') 
	`, id, id, id)
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
		AND actor_id = ?
		AND relationship = (SELECT id FROM ref_relationship WHERE type = 'blocked') 
	`, id)
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
		WHERE (user_id1 = ? OR user_id2 = ?)
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

func (c *Conn) GetContacts(id int) ([]entity.Friend, error) {
	rows, err := c.db.Query(`
		SELECT u.id, u.name 
		FROM user u 
		INNER JOIN (SELECT CASE 
			WHEN fr.user_id1 = 1 
				THEN fr.user_id2 
				ELSE fr.user_id1 
			END AS id 
		FROM friendship fr 
		WHERE 
			(fr.user_id1 = ? OR fr.user_id2 = ?) 
			AND fr.relationship = (SELECT id FROM ref_relationship WHERE type = "friend")) r 
		ON r.id = u.id
	`, id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.Friend
	for rows.Next() {
		var friend entity.Friend
		if err := rows.Scan(
			&friend.ID,
			&friend.Name,
		); err != nil {
			return nil, err
		}
		result = append(result, friend)
	}
	return result, rows.Err()
}
