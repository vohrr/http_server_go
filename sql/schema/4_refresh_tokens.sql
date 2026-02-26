-- +goose Up
CREATE TABLE refresh_tokens (
	token TEXT not null,
	created_at TIMESTAMP not null,
	updated_at TIMESTAMP not null,
	user_id uuid not null,
	expires_at TIMESTAMP not null,
	revoked_at TIMESTAMP null,
	CONSTRAINT fk_users
		FOREIGN KEY (user_id)
		REFERENCES users(id)
		ON DELETE CASCADE
);


-- +goose Down
DROP TABLE refresh_tokens;
