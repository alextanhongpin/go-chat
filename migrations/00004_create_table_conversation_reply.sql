-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS conversation_reply (
	id              int          AUTO_INCREMENT,
	text            varchar(140) NOT NULL,
	user_id         int,
	conversation_id int,
	PRIMARY KEY (id),
	FOREIGN KEY (user_id)         REFERENCES user(id),
	FOREIGN KEY (conversation_id) REFERENCES conversation(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4; 

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE conversation_reply;
