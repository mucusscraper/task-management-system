-- +goose Up
INSERT INTO users (id, name, role)
VALUES
    (1, 'Supervisor','SUPERVISOR'),
    (2, 'Worker 1', 'WORKER'),
    (3, 'Worker 2', 'WORKER');
SELECT setval(
    'users_id_seq',
    (SELECT MAX(id) FROM users)
);

-- +goose Down
TRUNCATE TABLE users RESTART IDENTITY CASCADE;
