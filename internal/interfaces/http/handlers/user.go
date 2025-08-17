package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"todo_list/config"
	"todo_list/internal/domain/models"
)

type UserRepository interface {
	GenerateJWTToken(user *models.UserWithCode, cfg *config.Config) (string, error)
	Login(ctx context.Context, user *models.UserWithCode, cfg *config.Config) (string, models.UserWithCode, error)
	Signup(ctx context.Context, user *models.User) error
	CreateUser(ctx context.Context, user *models.User) error
	GetUserInfo(ctx context.Context, id uint) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}

type UserHandler struct {
	db UserRepository
}

func NewUserHandler(db UserRepository) *UserHandler {
	return &UserHandler{db: db}
}

func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	var user models.UserWithCode

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	token, foundUser, err := uh.db.Login(r.Context(), &user, cfg)
	if err != nil {

		if err.Error() == "invalid password" {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if foundUser.Enable2FA {
		address := fmt.Sprintf("http://%s:%s/code", cfg.TwoFAHost, cfg.TwoFAPort)

		url := fmt.Sprintf("%s/%s", address, foundUser.TGUserName)

		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, "Failed to send 2FA request: ", http.StatusInternalServerError)
			return
		}

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Failed to send 2FA request: ", http.StatusInternalServerError)
			return
		}

		var resutl map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&resutl); err != nil {
			http.Error(w, "Failed to decode 2FA response: ", http.StatusInternalServerError)
			return
		}

		code := resutl["code"]
		if (code != user.Code) || (user.Code == "") {
			http.Error(w, "Invalid 2FA code", http.StatusUnauthorized)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (uh *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = uh.db.Signup(r.Context(), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (uh *UserHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := uh.db.GetUserInfo(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user.ID = userID

	err = uh.db.UpdateUser(r.Context(), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
