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

func (r *SQLUserRepository) GenerateJWTToken(user *models.UserWithCode, cfg *config.Config) (string, error) {

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

func (r *SQLUserRepository) Login(ctx context.Context, user *models.UserWithCode, cfg *config.Config) (token string, foundUser models.UserWithCode, err error) {

	query := "SELECT id, username, password, created_at, updated_at, email, enable_2fa, tg_username FROM users WHERE username = $1"

	err = r.db.QueryRow(ctx, query, user.Username).Scan(&foundUser.ID, &foundUser.Username, &foundUser.Password, &foundUser.CreatedAt, &foundUser.UpdatedAt, &foundUser.Email, &foundUser.Enable2FA, &foundUser.TGUserName)
	if err != nil {
		return "", foundUser, fmt.Errorf("failed to find user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		return "", foundUser, fmt.Errorf("invalid password: %w", err)
	}

	token, err = r.GenerateJWTToken(&foundUser, cfg)
	if err != nil {
		return "", foundUser, fmt.Errorf("could not generate token: %w", err)
	}

	return token, foundUser, nil
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
	query := "INSERT INTO users (username, password, created_at, updated_at, email, enable_2fa, tg_username) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	return r.db.QueryRow(ctx, query, user.Username, user.Password, user.CreatedAt, user.UpdatedAt, user.Email, user.Enable2FA, user.TGUserName).Scan(&user.ID)
}

func (r *SQLUserRepository) GetUserInfo(ctx context.Context, id uint) (*models.User, error) {
	query := "SELECT id, username, password, created_at, updated_at, email, enable_2fa, tg_username FROM users WHERE id = $1"
	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.Email, &user.Enable2FA, &user.TGUserName)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

func (r *SQLUserRepository) UpdateUser(ctx context.Context, user *models.User) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("could not hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	query := "UPDATE users SET username = $1, password = $2, enable_2fa = $3, tg_username = $4, updated_at = $5 WHERE id = $6"
	_, err = r.db.Exec(ctx, query, user.Username, user.Password, user.Enable2FA, user.TGUserName, time.Now(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
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
