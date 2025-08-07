-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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

INSERT INTO users (username, email, created_at, updated_at, password) VALUES
('user1', 'user1@example.com', NOW(), NOW(), 'password1'),
('user2', 'user2@example.com', NOW(), NOW(), 'password2'),
('user3', 'user3@example.com', NOW(), NOW(), 'password3');

INSERT INTO tasks (name, description, status, created_at, updated_at, user_id) VALUES
('Заметка 1', 'Описание заметки 1', 'pending', NOW(), NOW(), 1),
('Заметка 2', 'Описание заметки 2', 'pending', NOW(), NOW(), 1),
('Заметка 3', 'Описание заметки 3', 'pending', NOW(), NOW(), 2),
('Заметка 4', 'Описание заметки 4', 'pending', NOW(), NOW(), 3);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE tasks;
DROP TABLE users;
DELETE FROM tasks WHERE name IN ('Заметка 1', 'Заметка 2', 'Заметка 3', 'Заметка 4');
DELETE FROM users WHERE username IN ('user1', 'user2', 'user3');
-- +goose StatementEnd
