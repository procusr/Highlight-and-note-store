package services

import (
	"database/sql"
	"fmt"
	"mh-go/internal/models"
	"time"

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

// Register creates a new user account
func (s *AuthService) Register(email, password string) (*models.User, error) {
	// Check if user already exists
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		ID:       uuid.New(),
		Email:    email,
		Password: string(hashedPassword),
	}

	query := `
		INSERT INTO users (id, email, password)
		VALUES ($1, $2, $3)
		RETURNING created_at, updated_at
	`

	err = s.db.QueryRow(query, user.ID, user.Email, user.Password).
		Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Clear password from response
	user.Password = ""
	return user, nil
}

// Login authenticates a user
func (s *AuthService) Login(email, password string) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1`

	err := s.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid email or password")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Clear password from response
	user.Password = ""
	return &user, nil
}

// GenerateToken generates a JWT token for a user
func (s *AuthService) GenerateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (s *AuthService) ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid token claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid user_id claim")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user_id format: %w", err)
	}

	return userID, nil
}

// GetUserByID returns a user by ID
func (s *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, created_at, updated_at FROM users WHERE id = $1`

	err := s.db.QueryRow(query, userID).Scan(
		&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUserByEmail returns a user by email
func (s *AuthService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, created_at, updated_at FROM users WHERE email = $1`

	err := s.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// ResetPassword resets a user's password
func (s *AuthService) ResetPassword(token, password string) error {
	// Validate token
	userID, err := s.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("invalid or expired token")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	query := `UPDATE users SET password = $1 WHERE id = $2`
	_, err = s.db.Exec(query, string(hashedPassword), userID)
	if err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}

	return nil
}
