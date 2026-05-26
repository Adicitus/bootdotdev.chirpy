-- +goose Up
CREATE TABLE identities (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    auth TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE identities;
