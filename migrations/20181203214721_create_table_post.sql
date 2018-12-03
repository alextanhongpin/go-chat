-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS post (
	id int AUTO_INCREMENT,
	text text,
	user_id int,
	is_published bool NOT NULL DEFAULT false,
	created_at datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted_at datetime 	NOT NULL DEFAULT '1900-01-01 00:00:00',
	FOREIGN KEY (user_id) REFERENCES user (id),
	PRIMARY KEY (id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE post;
