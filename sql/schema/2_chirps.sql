-- +goose Up
CREATE TABLE chirps (
	id uuid not null primary key,
	created_at TIMESTAMP not null,
	updated_at TIMESTAMP not null,
	body VARCHAR(140) not null,
	user_id uuid not null,
	CONSTRAINT fk_users
		FOREIGN KEY (user_id)
		REFERENCES users (id)
		ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirps;
