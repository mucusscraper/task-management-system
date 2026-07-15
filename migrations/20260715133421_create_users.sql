-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('SUPERVISOR','WORKER'))
);

-- +goose Down
DROP TABLE users;
