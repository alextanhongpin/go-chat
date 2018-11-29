-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS user (
	id   int AUTO_INCREMENT,
	name varchar(65) NOT NULL,
	hashed_password varchar(255) NOT NULL,
	created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted_at datetime NOT NULL DEFAULT '1900-01-01 00:00:00', 
	email varchar(255) NOT NULL UNIQUE,
	PRIMARY KEY(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4; 

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE user;
