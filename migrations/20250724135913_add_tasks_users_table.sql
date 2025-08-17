-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    enable_2fa BOOLEAN DEFAULT FALSE,
    tg_username VARCHAR(255)
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    user_id INT REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO users (username, email, created_at, updated_at, password, enable_2fa, tg_username) VALUES
('user1', 'user1@example.com', NOW(), NOW(), '$2a$10$RWpNi.ka9VASMZpAEUygyuwLPMnA4/5u2NbmnKcuG5bNqHzFU5TPC', FALSE, 'user1_tg'),
('user2', 'user2@example.com', NOW(), NOW(), '$2a$10$RWpNi.ka9VASMZpAEUygyuwLPMnA4/5u2NbmnKcuG5bNqHzFU5TPC', FALSE, 'user2_tg'),
('user3', 'user3@example.com', NOW(), NOW(), '$2a$10$RWpNi.ka9VASMZpAEUygyuwLPMnA4/5u2NbmnKcuG5bNqHzFU5TPC', FALSE, 'user3_tg');

INSERT INTO tasks (name, description, status, created_at, updated_at, user_id) VALUES
('Заметка 1', 'Описание заметки 1', 'pending', NOW(), NOW(), 1),
('Заметка 2', 'Описание заметки 2', 'pending', NOW(), NOW(), 1),
('Заметка 3', 'Описание заметки 3', 'pending', NOW(), NOW(), 2),
('Заметка 4', 'Описание заметки 4', 'pending', NOW(), NOW(), 3);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
