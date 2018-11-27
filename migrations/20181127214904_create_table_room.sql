-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS room (
	id         int          AUTO_INCREMENT, 
	name       varchar(255),
	created_at datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted_at datetime 	NOT NULL DEFAULT '1900-01-01 00:00:00',
	PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4; 

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE room;
