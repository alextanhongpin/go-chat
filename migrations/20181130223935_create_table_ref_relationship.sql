-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS ref_relationship (
	id int PRIMARY KEY AUTO_INCREMENT,
	type varchar(255) NOT NULL
);

INSERT INTO ref_relationship (type) VALUES ('request');
INSERT INTO ref_relationship (type) VALUES ('friend');
INSERT INTO ref_relationship (type) VALUES ('block');

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE ref_relationship;
