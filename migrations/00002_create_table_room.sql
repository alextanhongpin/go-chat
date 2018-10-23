-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS room (
	id         int          AUTO_INCREMENT, 
	name       varchar(32),
	type       boolean      DEFAULT 0, -- 0 means 1-to-1, 1 means group.
	created_at datetime     DEFAULT CURRENT_TIMESTAMP,
	updated_at datetime     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted_at datetime,
	PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4; 

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE room;
