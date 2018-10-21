-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS user (
	id   int         AUTO_INCREMENT,
	name varchar(65) NOT NULL,
	-- hashed_password
	-- created_at
	-- updated_at
	-- deleted_at
	-- email
	PRIMARY KEY(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4; 

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE user;
