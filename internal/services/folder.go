package services

import (
	"database/sql"
	"fmt"
	"mh-go/internal/models"

	"github.com/google/uuid"
)

type FolderService struct {
	db *sql.DB
}

func NewFolderService(db *sql.DB) *FolderService {
	return &FolderService{db: db}
}

// GetRootFolders returns all root folders (folders without a parent) for a user
// Front Desk folder will always appear first
func (s *FolderService) GetRootFolders(userID uuid.UUID) ([]models.Folder, error) {
	// First ensure Front Desk folder exists
	if err := s.EnsureFrontDeskFolder(userID); err != nil {
		return nil, fmt.Errorf("failed to ensure Front Desk folder: %w", err)
	}

	query := `
		SELECT id, name, description, parent_id, user_id, folder_type, created_at, updated_at
		FROM folders
		WHERE user_id = $1 AND parent_id IS NULL
		ORDER BY
			CASE WHEN folder_type = 'front_desk' THEN 0 ELSE 1 END,
			name ASC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query root folders: %w", err)
	}
	defer rows.Close()

	var folders []models.Folder
	for rows.Next() {
		var folder models.Folder
		err := rows.Scan(
			&folder.ID, &folder.Name, &folder.Description,
			&folder.ParentID, &folder.UserID, &folder.FolderType, &folder.CreatedAt, &folder.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan folder: %w", err)
		}
		folders = append(folders, folder)
	}

	return folders, nil
}

// GetFoldersByParent returns all child folders for a given parent folder
func (s *FolderService) GetFoldersByParent(parentID, userID uuid.UUID) ([]models.Folder, error) {
	query := `
		SELECT id, name, description, parent_id, user_id, folder_type, created_at, updated_at
		FROM folders
		WHERE parent_id = $1 AND user_id = $2
		ORDER BY name ASC
	`

	rows, err := s.db.Query(query, parentID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query child folders: %w", err)
	}
	defer rows.Close()

	var folders []models.Folder
	for rows.Next() {
		var folder models.Folder
		err := rows.Scan(
			&folder.ID, &folder.Name, &folder.Description,
			&folder.ParentID, &folder.UserID, &folder.FolderType, &folder.CreatedAt, &folder.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan folder: %w", err)
		}
		folders = append(folders, folder)
	}

	return folders, nil
}

// GetFolderByID returns a folder by its ID
func (s *FolderService) GetFolderByID(folderID, userID uuid.UUID) (*models.Folder, error) {
	query := `
		SELECT id, name, description, parent_id, user_id, folder_type, created_at, updated_at
		FROM folders
		WHERE id = $1 AND user_id = $2
	`

	var folder models.Folder
	err := s.db.QueryRow(query, folderID, userID).Scan(
		&folder.ID, &folder.Name, &folder.Description,
		&folder.ParentID, &folder.UserID, &folder.FolderType, &folder.CreatedAt, &folder.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("folder not found")
		}
		return nil, fmt.Errorf("failed to get folder: %w", err)
	}

	return &folder, nil
}

// GetFolderWithPath returns a folder with its breadcrumb path
func (s *FolderService) GetFolderWithPath(folderID, userID uuid.UUID) (*models.FolderWithPath, error) {
	folder, err := s.GetFolderByID(folderID, userID)
	if err != nil {
		return nil, err
	}

	path, err := s.GetFolderPath(folderID, userID)
	if err != nil {
		return nil, err
	}

	return &models.FolderWithPath{
		Folder: folder,
		Path:   path,
	}, nil
}

// GetFolderPath returns the breadcrumb path for a folder
func (s *FolderService) GetFolderPath(folderID, userID uuid.UUID) ([]models.FolderPath, error) {
	query := `
		WITH RECURSIVE folder_path AS (
			SELECT id, name, parent_id, 1 as level
			FROM folders
			WHERE id = $1 AND user_id = $2

			UNION ALL

			SELECT f.id, f.name, f.parent_id, fp.level + 1
			FROM folders f
			INNER JOIN folder_path fp ON f.id = fp.parent_id
			WHERE f.user_id = $2
		)
		SELECT id, name FROM folder_path ORDER BY level DESC
	`

	rows, err := s.db.Query(query, folderID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get folder path: %w", err)
	}
	defer rows.Close()

	var path []models.FolderPath
	for rows.Next() {
		var item models.FolderPath
		err := rows.Scan(&item.ID, &item.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan path item: %w", err)
		}
		path = append(path, item)
	}

	return path, nil
}

// CreateFolder creates a new folder
func (s *FolderService) CreateFolder(folder *models.Folder) error {
	if folder.ID == uuid.Nil {
		folder.ID = uuid.New()
	}

	// Set default folder type if not specified
	if folder.FolderType == "" {
		folder.FolderType = "normal"
	}

	query := `
		INSERT INTO folders (id, name, description, parent_id, user_id, folder_type)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`

	err := s.db.QueryRow(query, folder.ID, folder.Name, folder.Description, folder.ParentID, folder.UserID, folder.FolderType).
		Scan(&folder.CreatedAt, &folder.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create folder: %w", err)
	}

	return nil
}

// UpdateFolder updates an existing folder
func (s *FolderService) UpdateFolder(folder *models.Folder) error {
	// Check if this folder can be edited
	if !folder.IsEditable() {
		return fmt.Errorf("this folder cannot be edited")
	}

	query := `
		UPDATE folders
		SET name = $2, description = $3, parent_id = $4
		WHERE id = $1 AND user_id = $5 AND folder_type != 'front_desk'
		RETURNING updated_at
	`

	err := s.db.QueryRow(query, folder.ID, folder.Name, folder.Description, folder.ParentID, folder.UserID).
		Scan(&folder.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update folder: %w", err)
	}

	return nil
}

// DeleteFolder deletes a folder and all its contents
func (s *FolderService) DeleteFolder(folderID, userID uuid.UUID) error {
	// First check if this is a Front Desk folder
	folder, err := s.GetFolderByID(folderID, userID)
	if err != nil {
		return fmt.Errorf("failed to get folder: %w", err)
	}

	if !folder.IsDeletable() {
		return fmt.Errorf("this folder cannot be deleted")
	}

	query := `DELETE FROM folders WHERE id = $1 AND user_id = $2`

	result, err := s.db.Exec(query, folderID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("folder not found or not authorized to delete")
	}

	return nil
}

// EnsureFrontDeskFolder creates the Front Desk folder if it doesn't exist
func (s *FolderService) EnsureFrontDeskFolder(userID uuid.UUID) error {
	// Check if Front Desk folder already exists
	checkQuery := `
		SELECT COUNT(*) FROM folders
		WHERE user_id = $1 AND folder_type = 'front_desk'
	`

	var count int
	err := s.db.QueryRow(checkQuery, userID).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check for Front Desk folder: %w", err)
	}

	// If it already exists, we're done
	if count > 0 {
		return nil
	}

	// Create the Front Desk folder
	frontDeskFolder := &models.Folder{
		ID:          uuid.New(),
		Name:        "Front Desk",
		Description: "Your main workspace for important items and quick access",
		ParentID:    nil,
		UserID:      userID,
		FolderType:  "front_desk",
	}

	return s.CreateFolder(frontDeskFolder)
}
