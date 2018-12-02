-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS user_room (
	user_id int,
	room_id int,
	created_at datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted_at datetime 	NOT NULL DEFAULT '1900-01-01 00:00:00',
	FOREIGN KEY (user_id) REFERENCES user(id),
	FOREIGN KEY (room_id) REFERENCES room(id),
	INDEX user_room_idx (user_id, room_id),
	CONSTRAINT uq_user_id_room_id UNIQUE (user_id, room_id),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4; 

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE user_room;
