-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS ref_relationship (
	status varchar(255) PRIMARY KEY 
);

INSERT INTO ref_relationship (status) VALUES ('request');
INSERT INTO ref_relationship (status) VALUES ('friend');
INSERT INTO ref_relationship (status) VALUES ('block');

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE ref_relationship;
