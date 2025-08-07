package router

import (
	"todo_list/internal/interfaces/http/handlers"
	"todo_list/internal/interfaces/http/repository"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewRouter(dataBase *pgxpool.Pool) *mux.Router {

	userRepository := repository.NewSQLUserRepository(dataBase)
	taskRepository := repository.NewTaskRepository(dataBase)

	userHandler := handlers.NewUserHandler(userRepository)
	taskHandler := handlers.NewTaskHandler(taskRepository)

	r := mux.NewRouter()

	r.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	r.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.GetTask).Methods("GET")
	r.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.DeleteTask).Methods("DELETE")
	r.HandleFunc("/users/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/users/signup", userHandler.Signup).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}", userHandler.GetUserByID).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", userHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id:[0-9]+}", userHandler.DeleteUser).Methods("DELETE")
	r.HandleFunc("/users", userHandler.ListUsers).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}/tasks", userHandler.ListTasksByUserID).Methods("GET")

	// Add more routes as needed

	return r
}
