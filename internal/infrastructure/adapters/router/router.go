package router

import (
	"net/http"
	"todo_list/config"
	"todo_list/internal/interfaces/http/handlers"
	"todo_list/internal/interfaces/http/repository"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewRouter(dataBase *pgxpool.Pool, cfg *config.Config) *mux.Router {

	userRepository := repository.NewSQLUserRepository(dataBase)
	taskRepository := repository.NewTaskRepository(dataBase)

	userHandler := handlers.NewUserHandler(userRepository)
	taskHandler := handlers.NewTaskHandler(taskRepository)

	r := mux.NewRouter()

	r.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	r.HandleFunc("/tasks", taskHandler.ListTasks).Methods("GET")
	r.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.GetTask).Methods("GET")
	r.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.DeleteTask).Methods("DELETE")
	r.HandleFunc("/users/login", func(w http.ResponseWriter, r *http.Request) {
		userHandler.Login(w, r, cfg)
	}).Methods("POST")
	r.HandleFunc("/users/signup", userHandler.Signup).Methods("POST")
	r.HandleFunc("/users", userHandler.GetUserInfo).Methods("GET")
	r.HandleFunc("/users", userHandler.UpdateUser).Methods("PUT")

	return r
}
