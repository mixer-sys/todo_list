package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"todo_list/internal/domain/models"

	"github.com/gorilla/mux"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *models.Task) error
	GetTask(ctx context.Context, id, userid uint) (*models.Task, error)
	UpdateTask(ctx context.Context, task *models.Task) error
	DeleteTask(ctx context.Context, id, userid uint) error
	ListTasks(ctx context.Context, userID uint) ([]models.Task, error)
}

type TaskHandler struct {
	db TaskRepository
}

func NewTaskHandler(db TaskRepository) *TaskHandler {
	return &TaskHandler{db: db}
}

func (th *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)

		return
	}

	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	userID, ok := r.Context().Value("userID").(uint)

	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)

		return
	}

	task.UserID = userID

	err := th.db.CreateTask(r.Context(), &task)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)

		return
	}
}

func (th *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]
	taskIDInt, err := strconv.Atoi(taskID)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("userID").(uint)

	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)

		return
	}

	task, err := th.db.GetTask(r.Context(), uint(taskIDInt), userID)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (th *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]
	taskIDInt, err := strconv.Atoi(taskID)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)

		return
	}

	task.UpdatedAt = time.Now()

	userID, ok := r.Context().Value("userID").(uint)

	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)

		return
	}

	task.UserID = userID
	task.ID = uint(taskIDInt)

	err = th.db.UpdateTask(r.Context(), &task)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (th *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]
	taskIDInt, err := strconv.Atoi(taskID)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("userID").(uint)

	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)

		return
	}

	err = th.db.DeleteTask(r.Context(), uint(taskIDInt), userID)
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (th *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tasks, err := th.db.ListTasks(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
