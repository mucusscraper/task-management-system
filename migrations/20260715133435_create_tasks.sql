-- +goose Up
CREATE TABLE tasks(
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL CHECK ( status IN ('CREATED','ASSIGNED','IN_PROGRESS','COMPLETED')),
    assigned_to INTEGER NULL REFERENCES users(id),
    created_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE tasks CASCADE;
