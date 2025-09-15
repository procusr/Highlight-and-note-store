package services

import (
	"database/sql"
	"fmt"
	"mh-go/internal/models"

	"github.com/google/uuid"
)

type ContentService struct {
	db *sql.DB
}

func NewContentService(db *sql.DB) *ContentService {
	return &ContentService{db: db}
}

// GetContentByFolderID returns all content for a specific folder
func (s *ContentService) GetContentByFolderID(folderID, userID uuid.UUID) ([]models.Content, error) {
    query := `
		SELECT id, title, description, question, body, content_type, folder_id, user_id, created_at, updated_at
		FROM content
		WHERE folder_id = $1 AND user_id = $2
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, folderID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query content: %w", err)
	}
	defer rows.Close()

	var contents []models.Content
	for rows.Next() {
		var content models.Content
		err := rows.Scan(
			&content.ID, &content.Title, &content.Description, &content.Question, &content.Body,
			&content.ContentType, &content.FolderID, &content.UserID,
			&content.CreatedAt, &content.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan content: %w", err)
		}
		contents = append(contents, content)
	}

	return contents, nil
}

// GetContentByID returns content by its ID
func (s *ContentService) GetContentByID(contentID, userID uuid.UUID) (*models.Content, error) {
    query := `
		SELECT id, title, description, question, body, content_type, folder_id, user_id, created_at, updated_at
		FROM content
		WHERE id = $1 AND user_id = $2
	`

	var content models.Content
	err := s.db.QueryRow(query, contentID, userID).Scan(
		&content.ID, &content.Title, &content.Description, &content.Question, &content.Body,
		&content.ContentType, &content.FolderID, &content.UserID,
		&content.CreatedAt, &content.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("content not found")
		}
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	return &content, nil
}

// CreateContent creates new content
func (s *ContentService) CreateContent(content *models.Content) error {
    content.ID = uuid.New()

    query := `
		INSERT INTO content (id, title, description, question, body, content_type, folder_id, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at
	`

    err := s.db.QueryRow(query,
        content.ID, content.Title, content.Description, content.Question, content.Body,
        content.ContentType, content.FolderID, content.UserID).
        Scan(&content.CreatedAt, &content.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create content: %w", err)
	}

	return nil
}

// UpdateContent updates existing content
func (s *ContentService) UpdateContent(content *models.Content) error {
    query := `
		UPDATE content
		SET title = $2, description = $3, question = $4, body = $5, content_type = $6
		WHERE id = $1 AND user_id = $7
		RETURNING updated_at
	`

    err := s.db.QueryRow(query,
        content.ID, content.Title, content.Description, content.Question, content.Body,
        content.ContentType, content.UserID).
        Scan(&content.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update content: %w", err)
	}

	return nil
}

// DeleteContent deletes content by ID
func (s *ContentService) DeleteContent(contentID, userID uuid.UUID) error {
	query := `DELETE FROM content WHERE id = $1 AND user_id = $2`

	result, err := s.db.Exec(query, contentID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete content: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("content not found or not authorized to delete")
	}

	return nil
}

// SearchContent searches for content by title or description
func (s *ContentService) SearchContent(query string, userID uuid.UUID) ([]models.Content, error) {
    searchQuery := `
		SELECT id, title, description, question, body, content_type, folder_id, user_id, created_at, updated_at
		FROM content
		WHERE user_id = $1 AND (title ILIKE $2 OR description ILIKE $2 OR body ILIKE $2 OR question ILIKE $2)
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(searchQuery, userID, "%"+query+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search content: %w", err)
	}
	defer rows.Close()

	var contents []models.Content
	for rows.Next() {
		var content models.Content
		err := rows.Scan(
			&content.ID, &content.Title, &content.Description, &content.Question, &content.Body,
			&content.ContentType, &content.FolderID, &content.UserID,
			&content.CreatedAt, &content.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan content: %w", err)
		}
		contents = append(contents, content)
	}

	return contents, nil
}
