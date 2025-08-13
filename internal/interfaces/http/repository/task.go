package repository

import (
	"context"
	"fmt"
	"time"
	"todo_list/internal/domain/models"
	"todo_list/internal/interfaces/http/handlers"

	"github.com/jackc/pgx/v4/pgxpool"
)

type SQLTaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) handlers.TaskRepository {
	return &SQLTaskRepository{db: db}
}

func (r *SQLTaskRepository) CreateTask(ctx context.Context, task *models.Task) error {
	query := "INSERT INTO tasks (name, description, status, created_at, updated_at, user_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	return r.db.QueryRow(ctx, query, task.Name, task.Description, task.Status, task.CreatedAt, task.UpdatedAt, task.UserID).Scan(&task.ID)
}

func (r *SQLTaskRepository) GetTask(ctx context.Context, id, userID uint) (*models.Task, error) {
	query := "SELECT id, name, description, status, created_at, updated_at, user_id FROM tasks WHERE id = $1 AND user_id = $2"
	var task models.Task
	err := r.db.QueryRow(ctx, query, id, userID).Scan(&task.ID, &task.Name, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt, &task.UserID)

	if err != nil {

		return nil, fmt.Errorf("failed to get task by ID: %w", err)
	}

	return &task, nil
}

func (r *SQLTaskRepository) UpdateTask(ctx context.Context, task *models.Task) error {
	query := "UPDATE tasks SET name = $1, description = $2, status = $3, updated_at = $4 WHERE id = $5 AND user_id = $6"
	_, err := r.db.Exec(ctx, query, task.Name, task.Description, task.Status, time.Now(), task.ID, task.UserID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}
	return nil
}

func (r *SQLTaskRepository) DeleteTask(ctx context.Context, id, userID uint) error {
	query := "DELETE FROM tasks WHERE id = $1 AND user_id = $2"
	_, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

func (r *SQLTaskRepository) ListTasks(ctx context.Context, userID uint) ([]models.Task, error) {
	query := "SELECT id, name, description, status, created_at, updated_at, user_id FROM tasks WHERE user_id = $1"
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks by user ID: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt, &task.UserID); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
