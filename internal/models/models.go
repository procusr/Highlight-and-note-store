package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password_hash"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserFolder struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	IsDefault   bool      `json:"is_default" db:"is_default"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type UserContent struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	FolderID     *uuid.UUID `json:"folder_id" db:"folder_id"`
	Title        *string    `json:"title" db:"title"`
	Content      string     `json:"content" db:"content"`
	Anchor       *string    `json:"anchor" db:"anchor"`
	IsArchived   bool       `json:"is_archived" db:"is_archived"`
	IsFavorite   bool       `json:"is_favorite" db:"is_favorite"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	SearchVector *string    `json:"-" db:"search_vector"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token           string `json:"token" validate:"required"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type CreateFolderRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description"`
}

type UpdateFolderRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description"`
}

type CreateContentRequest struct {
	FolderID *uuid.UUID `json:"folder_id"`
	Title    *string    `json:"title"`
	Content  string     `json:"content" validate:"required"`
	Anchor   *string    `json:"anchor"`
}

type UpdateContentRequest struct {
	FolderID *uuid.UUID `json:"folder_id"`
	Title    *string    `json:"title"`
	Content  string     `json:"content" validate:"required"`
	Anchor   *string    `json:"anchor"`
}
