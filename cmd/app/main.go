package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"todo_list/internal/models"
)

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("/data/todo.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database" + err.Error())
	}

	if err := db.AutoMigrate(&models.Task{}); err != nil {
		panic("failed to migrate database")
	}

	r := mux.NewRouter()
	// Define your routes here
	r.HandleFunc("/tasks", createTaskHandler).Methods("POST")
	r.HandleFunc("/tasks/{id}", getTaskHandler).Methods("GET")
	r.HandleFunc("/tasks", listTasksHandler).Methods("GET")
	r.HandleFunc("/tasks/{id}", updateTaskHandler).Methods("PUT")
	r.HandleFunc("/tasks/{id}", deleteTaskHandler).Methods("DELETE")
	http.ListenAndServe(":8080", r)
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Handler logic to create a task
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := db.Create(&task).Error; err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Handler logic to get a task by ID
	vars := mux.Vars(r)
	var task models.Task
	if err := db.First(&task, vars["id"]).Error; err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)

}
func listTasksHandler(w http.ResponseWriter, r *http.Request) {
	// Handler logic to list all tasks
	var tasks []models.Task
	if err := db.Find(&tasks).Error; err != nil {
		http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Handler logic to update a task by ID
	vars := mux.Vars(r)
	var task models.Task
	if err := db.First(&task, vars["id"]).Error; err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := db.Save(&task).Error; err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Handler logic to delete a task by ID
	vars := mux.Vars(r)
	var task models.Task
	if err := db.First(&task, vars["id"]).Error; err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	if err := db.Delete(&task).Error; err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted successfully"})
}
