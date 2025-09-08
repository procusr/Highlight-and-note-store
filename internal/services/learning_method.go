package services

import (
	"database/sql"
	"fmt"
	"mh-go/internal/models"
	"strings"

	"github.com/google/uuid"
)

type LearningMethodService struct {
	db *sql.DB
}

func NewLearningMethodService(db *sql.DB) *LearningMethodService {
	return &LearningMethodService{db: db}
}

// GetAllLearningMethods returns all available learning methods
func (s *LearningMethodService) GetAllLearningMethods() ([]models.LearningMethod, error) {
	query := `
		SELECT id, name, description, icon, method_key, ai_prompt_template, created_at, updated_at
		FROM learning_methods
		ORDER BY name
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query learning methods: %w", err)
	}
	defer rows.Close()

	var methods []models.LearningMethod
	for rows.Next() {
		var method models.LearningMethod
		err := rows.Scan(
			&method.ID, &method.Name, &method.Description, &method.Icon,
			&method.MethodKey, &method.AiPromptTemplate, &method.CreatedAt, &method.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan learning method: %w", err)
		}
		methods = append(methods, method)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate learning methods: %w", err)
	}

	return methods, nil
}

// GetLearningMethodByID returns a learning method by ID
func (s *LearningMethodService) GetLearningMethodByID(methodID uuid.UUID) (*models.LearningMethod, error) {
	var method models.LearningMethod
	query := `
		SELECT id, name, description, icon, method_key, ai_prompt_template, created_at, updated_at
		FROM learning_methods
		WHERE id = $1
	`

	err := s.db.QueryRow(query, methodID).Scan(
		&method.ID, &method.Name, &method.Description, &method.Icon,
		&method.MethodKey, &method.AiPromptTemplate, &method.CreatedAt, &method.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("learning method not found")
		}
		return nil, fmt.Errorf("failed to get learning method: %w", err)
	}

	return &method, nil
}

// GetLearningMethodByKey returns a learning method by method key
func (s *LearningMethodService) GetLearningMethodByKey(methodKey string) (*models.LearningMethod, error) {
	var method models.LearningMethod
	query := `
		SELECT id, name, description, icon, method_key, ai_prompt_template, created_at, updated_at
		FROM learning_methods
		WHERE method_key = $1
	`

	err := s.db.QueryRow(query, methodKey).Scan(
		&method.ID, &method.Name, &method.Description, &method.Icon,
		&method.MethodKey, &method.AiPromptTemplate, &method.CreatedAt, &method.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("learning method not found")
		}
		return nil, fmt.Errorf("failed to get learning method: %w", err)
	}

	return &method, nil
}

// CreateLearningMethod creates a new learning method
func (s *LearningMethodService) CreateLearningMethod(method *models.LearningMethod) error {
	method.ID = uuid.New()

	query := `
		INSERT INTO learning_methods (id, name, description, icon, method_key, ai_prompt_template)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`

	err := s.db.QueryRow(
		query, method.ID, method.Name, method.Description, 
		method.Icon, method.MethodKey, method.AiPromptTemplate,
	).Scan(&method.CreatedAt, &method.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create learning method: %w", err)
	}

	return nil
}

// UpdateLearningMethod updates an existing learning method
func (s *LearningMethodService) UpdateLearningMethod(method *models.LearningMethod) error {
	query := `
		UPDATE learning_methods
		SET name = $2, description = $3, icon = $4, method_key = $5, 
		    ai_prompt_template = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`

	err := s.db.QueryRow(
		query, method.ID, method.Name, method.Description, 
		method.Icon, method.MethodKey, method.AiPromptTemplate,
	).Scan(&method.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("learning method not found")
		}
		return fmt.Errorf("failed to update learning method: %w", err)
	}

	return nil
}

// DeleteLearningMethod deletes a learning method
func (s *LearningMethodService) DeleteLearningMethod(methodID uuid.UUID) error {
	query := `DELETE FROM learning_methods WHERE id = $1`

	result, err := s.db.Exec(query, methodID)
	if err != nil {
		return fmt.Errorf("failed to delete learning method: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("learning method not found")
	}

	return nil
}

// GenerateAIPrompt generates the AI prompt by replacing placeholders
func (s *LearningMethodService) GenerateAIPrompt(methodID uuid.UUID, topic string) (string, error) {
	method, err := s.GetLearningMethodByID(methodID)
	if err != nil {
		return "", fmt.Errorf("failed to get learning method: %w", err)
	}

	// Replace {TOPIC} placeholder with the actual topic
	prompt := strings.ReplaceAll(method.AiPromptTemplate, "{TOPIC}", topic)
	
	return prompt, nil
}

// GenerateAIPromptByKey generates the AI prompt by method key
func (s *LearningMethodService) GenerateAIPromptByKey(methodKey string, topic string) (string, error) {
	method, err := s.GetLearningMethodByKey(methodKey)
	if err != nil {
		return "", fmt.Errorf("failed to get learning method: %w", err)
	}

	// Replace {TOPIC} placeholder with the actual topic
	prompt := strings.ReplaceAll(method.AiPromptTemplate, "{TOPIC}", topic)
	
	return prompt, nil
}