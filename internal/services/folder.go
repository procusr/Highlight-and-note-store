package services

import (
	"database/sql"
	"fmt"
	"time"

	"mh-go/internal/models"

	"github.com/google/uuid"
)

type FolderService struct {
	db *sql.DB
}

func NewFolderService(db *sql.DB) *FolderService {
	return &FolderService{db: db}
}

func (s *FolderService) GetByUserID(userID uuid.UUID) ([]*models.UserFolder, error) {
	rows, err := s.db.Query(`
		SELECT id, user_id, name, description, is_default, created_at, updated_at
		FROM public.user_folders
		WHERE user_id = $1
		ORDER BY is_default DESC, name ASC`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get folders: %w", err)
	}
	defer rows.Close()

	var folders []*models.UserFolder
	for rows.Next() {
		folder := &models.UserFolder{}
		err := rows.Scan(&folder.ID, &folder.UserID, &folder.Name, &folder.Description,
			&folder.IsDefault, &folder.CreatedAt, &folder.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan folder: %w", err)
		}
		folders = append(folders, folder)
	}

	return folders, nil
}

func (s *FolderService) Create(userID uuid.UUID, req *models.CreateFolderRequest) (*models.UserFolder, error) {
	folder := &models.UserFolder{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		IsDefault:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := s.db.Exec(`
		INSERT INTO public.user_folders (id, user_id, name, description, is_default, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		folder.ID, folder.UserID, folder.Name, folder.Description, folder.IsDefault,
		folder.CreatedAt, folder.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create folder: %w", err)
	}

	return folder, nil
}

func (s *FolderService) GetByID(id, userID uuid.UUID) (*models.UserFolder, error) {
	folder := &models.UserFolder{}
	err := s.db.QueryRow(`
		SELECT id, user_id, name, description, is_default, created_at, updated_at
		FROM public.user_folders
		WHERE id = $1 AND user_id = $2`, id, userID).Scan(
		&folder.ID, &folder.UserID, &folder.Name, &folder.Description,
		&folder.IsDefault, &folder.CreatedAt, &folder.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("folder not found: %w", err)
	}
	return folder, nil
}

func (s *FolderService) Update(id, userID uuid.UUID, req *models.UpdateFolderRequest) (*models.UserFolder, error) {
	_, err := s.db.Exec(`
		UPDATE public.user_folders
		SET name = $1, description = $2, updated_at = $3
		WHERE id = $4 AND user_id = $5`,
		req.Name, req.Description, time.Now(), id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to update folder: %w", err)
	}

	return s.GetByID(id, userID)
}

func (s *FolderService) Delete(id, userID uuid.UUID) error {
	// Check if it's the default folder
	var isDefault bool
	err := s.db.QueryRow("SELECT is_default FROM public.user_folders WHERE id = $1 AND user_id = $2",
		id, userID).Scan(&isDefault)
	if err != nil {
		return fmt.Errorf("folder not found: %w", err)
	}

	if isDefault {
		return fmt.Errorf("cannot delete default folder")
	}

	// Move content to default folder
	_, err = s.db.Exec(`
		UPDATE public.user_content
		SET folder_id = (
			SELECT id FROM public.user_folders
			WHERE user_id = $1 AND is_default = true
		)
		WHERE folder_id = $2`, userID, id)
	if err != nil {
		return fmt.Errorf("failed to move content: %w", err)
	}

	// Delete folder
	_, err = s.db.Exec("DELETE FROM public.user_folders WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}

	return nil
}
