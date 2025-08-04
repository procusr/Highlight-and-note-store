package services

import (
	"database/sql"
	"fmt"
	"time"

	"mh-go/internal/models"

	"github.com/google/uuid"
)

type ContentService struct {
	db *sql.DB
}

func NewContentService(db *sql.DB) *ContentService {
	return &ContentService{db: db}
}

func (s *ContentService) GetByUserID(userID uuid.UUID, folderID *uuid.UUID) ([]*models.UserContent, error) {
	var query string
	var args []interface{}

	if folderID != nil {
		query = `
			SELECT id, user_id, folder_id, title, content, anchor, is_archived, is_favorite, created_at, updated_at
			FROM public.user_content
			WHERE user_id = $1 AND folder_id = $2 AND is_archived = false
			ORDER BY is_favorite DESC, created_at DESC`
		args = []interface{}{userID, *folderID}
	} else {
		query = `
			SELECT id, user_id, folder_id, title, content, anchor, is_archived, is_favorite, created_at, updated_at
			FROM public.user_content
			WHERE user_id = $1 AND is_archived = false
			ORDER BY is_favorite DESC, created_at DESC`
		args = []interface{}{userID}
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get content: %w", err)
	}
	defer rows.Close()

	var contents []*models.UserContent
	for rows.Next() {
		content := &models.UserContent{}
		err := rows.Scan(&content.ID, &content.UserID, &content.FolderID, &content.Title,
			&content.Content, &content.Anchor, &content.IsArchived, &content.IsFavorite,
			&content.CreatedAt, &content.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan content: %w", err)
		}
		contents = append(contents, content)
	}

	return contents, nil
}

func (s *ContentService) Create(userID uuid.UUID, req *models.CreateContentRequest) (*models.UserContent, error) {
	content := &models.UserContent{
		ID:         uuid.New(),
		UserID:     userID,
		FolderID:   req.FolderID,
		Title:      req.Title,
		Content:    req.Content,
		Anchor:     req.Anchor,
		IsArchived: false,
		IsFavorite: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err := s.db.Exec(`
		INSERT INTO public.user_content (id, user_id, folder_id, title, content, anchor, is_archived, is_favorite, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		content.ID, content.UserID, content.FolderID, content.Title, content.Content,
		content.Anchor, content.IsArchived, content.IsFavorite, content.CreatedAt, content.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create content: %w", err)
	}

	return content, nil
}

func (s *ContentService) GetByID(id, userID uuid.UUID) (*models.UserContent, error) {
	content := &models.UserContent{}
	err := s.db.QueryRow(`
		SELECT id, user_id, folder_id, title, content, anchor, is_archived, is_favorite, created_at, updated_at
		FROM public.user_content
		WHERE id = $1 AND user_id = $2`, id, userID).Scan(
		&content.ID, &content.UserID, &content.FolderID, &content.Title, &content.Content,
		&content.Anchor, &content.IsArchived, &content.IsFavorite, &content.CreatedAt, &content.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("content not found: %w", err)
	}
	return content, nil
}

func (s *ContentService) Update(id, userID uuid.UUID, req *models.UpdateContentRequest) (*models.UserContent, error) {
	_, err := s.db.Exec(`
		UPDATE public.user_content
		SET folder_id = $1, title = $2, content = $3, anchor = $4, updated_at = $5
		WHERE id = $6 AND user_id = $7`,
		req.FolderID, req.Title, req.Content, req.Anchor, time.Now(), id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to update content: %w", err)
	}

	return s.GetByID(id, userID)
}

func (s *ContentService) Delete(id, userID uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM public.user_content WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete content: %w", err)
	}
	return nil
}

func (s *ContentService) Archive(id, userID uuid.UUID) error {
	_, err := s.db.Exec(`
		UPDATE public.user_content
		SET is_archived = true, updated_at = $1
		WHERE id = $2 AND user_id = $3`, time.Now(), id, userID)
	if err != nil {
		return fmt.Errorf("failed to archive content: %w", err)
	}
	return nil
}

func (s *ContentService) ToggleFavorite(id, userID uuid.UUID) error {
	_, err := s.db.Exec(`
		UPDATE public.user_content
		SET is_favorite = NOT is_favorite, updated_at = $1
		WHERE id = $2 AND user_id = $3`, time.Now(), id, userID)
	if err != nil {
		return fmt.Errorf("failed to toggle favorite: %w", err)
	}
	return nil
}
