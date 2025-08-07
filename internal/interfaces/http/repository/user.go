package repository

import (
	"context"
	"fmt"
	"time"
	"todo_list/config"
	"todo_list/internal/domain/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type SQLUserRepository struct {
	db *pgxpool.Pool
}

func NewSQLUserRepository(db *pgxpool.Pool) *SQLUserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) GenerateJWTToken(user *models.User, cfg *config.Config) (string, error) {

	expirationTime := time.Now().Add(time.Hour * time.Duration(cfg.ExpirationTimeHours))

	claims := &models.Claims{
		ID:       user.ID,
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cfg.JWTSecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (r *SQLUserRepository) Login(ctx context.Context, user *models.User, cfg *config.Config) (string, error) {

	query := "SELECT id, username, password, created_at, updated_at, email FROM users WHERE username = $1"
	var foundUser models.User
	err := r.db.QueryRow(ctx, query, user.Username).Scan(&foundUser.ID, &foundUser.Username, &foundUser.Password, &foundUser.CreatedAt, &foundUser.UpdatedAt, &foundUser.Email)
	if err != nil {
		return "", fmt.Errorf("failed to find user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		return "", fmt.Errorf("invalid password: %w", err)
	}

	token, err := r.GenerateJWTToken(&foundUser, cfg)
	if err != nil {
		return "", fmt.Errorf("could not generate token: %w", err)
	}

	return token, nil
}

func (r *SQLUserRepository) Signup(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("could not hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	return r.CreateUser(ctx, user)
}

func (r *SQLUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users (username, password, created_at, updated_at, email) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	return r.db.QueryRow(ctx, query, user.Username, user.Password, user.CreatedAt, user.UpdatedAt, user.Email).Scan(&user.ID)
}

func (r *SQLUserRepository) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	query := "SELECT id, username, password, created_at, updated_at, email FROM users WHERE id = $1"
	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

func (r *SQLUserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := "UPDATE users SET username = $1, password = $2, updated_at = $3 WHERE id = $4"
	_, err := r.db.Exec(ctx, query, user.Username, user.Password, time.Now(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *SQLUserRepository) DeleteUser(ctx context.Context, id uint) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (r *SQLUserRepository) ListUsers(ctx context.Context) ([]models.User, error) {
	query := "SELECT id, username, password, created_at, updated_at, email FROM users"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.Email); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *SQLUserRepository) ListTasksByUserID(ctx context.Context, userID uint) ([]models.Task, error) {
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

func ExtractClaims(tokenStr string, cfg *config.Config) (*models.Claims, error) {
	hmacSecretString := cfg.JWTSecretKey
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parse token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &models.Claims{
			ID:       uint(claims["id"].(float64)),
			UserID:   uint(claims["user_id"].(float64)),
			Username: claims["username"].(string),
			Email:    claims["email"].(string),
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: int64(claims["exp"].(float64)),
				IssuedAt:  int64(claims["iat"].(float64)),
			},
		}, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
