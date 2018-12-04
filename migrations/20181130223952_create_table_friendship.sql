-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS friendship (
	user_id1 int,
	user_id2 int,
	actor_id int,
	relationship varchar(32),

	created_at datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

	CONSTRAINT check_one_way CHECK (user_id1 < user_id2),
	CONSTRAINT uq_user_id_1_user_id_2 UNIQUE (user_id1, user_id2),
	FOREIGN KEY (user_id1) REFERENCES user (id) ON UPDATE CASCADE ON DELETE CASCADE,
	FOREIGN KEY (user_id2) REFERENCES user (id) ON UPDATE CASCADE ON DELETE CASCADE,
	FOREIGN KEY (relationship) REFERENCES ref_relationship (status) ON UPDATE CASCADE ON DELETE SET NULL,
	PRIMARY KEY (user_id1, user_id2)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE friendship;
