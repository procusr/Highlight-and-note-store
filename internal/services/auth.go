package services

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"mh-go/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db        *sql.DB
	jwtSecret []byte
}

func NewAuthService(db *sql.DB, jwtSecret string) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *AuthService) Register(email, password string) (*models.User, error) {
	// Check if user already exists
	var existingID uuid.UUID
	err := s.db.QueryRow("SELECT id FROM auth.users WHERE email = $1", email).Scan(&existingID)
	if err == nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		ID:        uuid.New(),
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.db.Exec(`
		INSERT INTO auth.users (id, email, encrypted_password, email_confirmed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, user.Email, user.Password, time.Now(), user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create default folder
	_, err = s.db.Exec(`
		INSERT INTO public.user_folders (id, user_id, name, description, is_default, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		uuid.New(), user.ID, "Default", "Default folder for your content", true, time.Now(), time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to create default folder: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	var user models.User
	err := s.db.QueryRow(`
		SELECT id, email, encrypted_password, created_at, updated_at
		FROM auth.users
		WHERE email = $1`, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 1 week
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return uuid.Nil, fmt.Errorf("invalid user_id in token")
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return uuid.Nil, fmt.Errorf("invalid user_id format: %w", err)
		}

		return userID, nil
	}

	return uuid.Nil, fmt.Errorf("invalid token")
}

func (s *AuthService) GenerateResetToken(email string) (string, error) {
	// Check if user exists
	var userID uuid.UUID
	err := s.db.QueryRow("SELECT id FROM auth.users WHERE email = $1", email).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	// Generate random token
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	token := hex.EncodeToString(bytes)

	// Store token in database (you would need to create a password_reset_tokens table)
	// For now, we'll just return the token
	return token, nil
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
	// Verify token and get user ID (implementation depends on your reset token storage)
	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password (this is a placeholder - you'd need to implement token verification)
	// _, err = s.db.Exec("UPDATE auth.users SET encrypted_password = $1 WHERE id = $2", hashedPassword, userID)
	_ = hashedPassword // placeholder to avoid unused variable error

	return fmt.Errorf("password reset not fully implemented yet")
}
