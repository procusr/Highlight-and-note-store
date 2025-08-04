package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Folder represents a folder in the system
type Folder struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	ParentID    *uuid.UUID `json:"parent_id" db:"parent_id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`

	// Relationships
	Parent   *Folder   `json:"parent,omitempty"`
	Children []Folder  `json:"children,omitempty"`
	Contents []Content `json:"contents,omitempty"`
	User     *User     `json:"user,omitempty"`
}

// Content represents content within a folder
type Content struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Body        string    `json:"body" db:"body"`
	ContentType string    `json:"content_type" db:"content_type"` // text, image, file, etc.
	FolderID    uuid.UUID `json:"folder_id" db:"folder_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Relationships
	Folder *Folder `json:"folder,omitempty"`
	User   *User   `json:"user,omitempty"`
}

// FolderPath represents a breadcrumb path item
type FolderPath struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// FolderWithPath includes breadcrumb navigation
type FolderWithPath struct {
	*Folder
	Path []FolderPath `json:"path"`
}
