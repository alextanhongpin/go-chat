-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS room (
	id         int          AUTO_INCREMENT, 
	name       varchar(32),
	created_at datetime     DEFAULT UTC_TIMESTAMP,
	updated_at datetime     DEFAULT UTC_TIMESTAMP ON UPDATE UTC_TIMESTAMP, 
	deleted_at datetime,
	PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4; 

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE room;
