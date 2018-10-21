-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS room (
	id      int AUTO_INCREMENT, 
	-- createdBy
	PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4; 

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE room;
