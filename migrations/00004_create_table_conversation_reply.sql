-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS conversation (
	id         int          AUTO_INCREMENT,
	text       varchar(140) NOT NULL,
	user_id    int,
	room_id    int,
	created_at datetime     DEFAULT UTC_TIMESTAMP,
	updated_at datetime     DEFAULT UTC_TIMESTAMP ON UPDATE UTC_TIMESTAMP, 
	deleted_at datetime,
	PRIMARY KEY (id),
	FOREIGN KEY (user_id) REFERENCES user(id),
	FOREIGN KEY (room_id) REFERENCES room(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4; 

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE conversation_reply;
