-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS post_reaction (
	post_id int,
	user_id int,
	--  vote bit NOT NULL DEFAULT 0,
	reaction tinyint NOT NULL DEFAULT -1,
	CONSTRAINT UNIQUE (user_id, post_id),
	FOREIGN KEY (post_id) REFERENCES post (id) ON UPDATE CASCADE ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES user (id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE post_reaction;
