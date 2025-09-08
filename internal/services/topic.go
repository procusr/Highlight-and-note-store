package services

import (
	"database/sql"
	"fmt"
	"mh-go/internal/models"

	"github.com/google/uuid"
)

type TopicService struct {
	db *sql.DB
}

func NewTopicService(db *sql.DB) *TopicService {
	return &TopicService{db: db}
}

// GetAllTopics returns all available topics
func (s *TopicService) GetAllTopics() ([]models.Topic, error) {
	query := `
		SELECT id, name, description, category, difficulty_level, created_at, updated_at
		FROM topics
		ORDER BY category, name
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query topics: %w", err)
	}
	defer rows.Close()

	var topics []models.Topic
	for rows.Next() {
		var topic models.Topic
		err := rows.Scan(
			&topic.ID, &topic.Name, &topic.Description, &topic.Category,
			&topic.DifficultyLevel, &topic.CreatedAt, &topic.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan topic: %w", err)
		}
		topics = append(topics, topic)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate topics: %w", err)
	}

	return topics, nil
}

// GetTopicByID returns a topic by ID
func (s *TopicService) GetTopicByID(topicID uuid.UUID) (*models.Topic, error) {
	var topic models.Topic
	query := `
		SELECT id, name, description, category, difficulty_level, created_at, updated_at
		FROM topics
		WHERE id = $1
	`

	err := s.db.QueryRow(query, topicID).Scan(
		&topic.ID, &topic.Name, &topic.Description, &topic.Category,
		&topic.DifficultyLevel, &topic.CreatedAt, &topic.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("topic not found")
		}
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}

	return &topic, nil
}

// GetTopicsByCategory returns topics filtered by category
func (s *TopicService) GetTopicsByCategory(category string) ([]models.Topic, error) {
	query := `
		SELECT id, name, description, category, difficulty_level, created_at, updated_at
		FROM topics
		WHERE category = $1
		ORDER BY difficulty_level, name
	`

	rows, err := s.db.Query(query, category)
	if err != nil {
		return nil, fmt.Errorf("failed to query topics by category: %w", err)
	}
	defer rows.Close()

	var topics []models.Topic
	for rows.Next() {
		var topic models.Topic
		err := rows.Scan(
			&topic.ID, &topic.Name, &topic.Description, &topic.Category,
			&topic.DifficultyLevel, &topic.CreatedAt, &topic.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan topic: %w", err)
		}
		topics = append(topics, topic)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate topics: %w", err)
	}

	return topics, nil
}

// CreateTopic creates a new topic
func (s *TopicService) CreateTopic(topic *models.Topic) error {
	topic.ID = uuid.New()

	query := `
		INSERT INTO topics (id, name, description, category, difficulty_level)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at
	`

	err := s.db.QueryRow(
		query, topic.ID, topic.Name, topic.Description, 
		topic.Category, topic.DifficultyLevel,
	).Scan(&topic.CreatedAt, &topic.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}

	return nil
}

// UpdateTopic updates an existing topic
func (s *TopicService) UpdateTopic(topic *models.Topic) error {
	query := `
		UPDATE topics
		SET name = $2, description = $3, category = $4, difficulty_level = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`

	err := s.db.QueryRow(
		query, topic.ID, topic.Name, topic.Description, 
		topic.Category, topic.DifficultyLevel,
	).Scan(&topic.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("topic not found")
		}
		return fmt.Errorf("failed to update topic: %w", err)
	}

	return nil
}

// DeleteTopic deletes a topic
func (s *TopicService) DeleteTopic(topicID uuid.UUID) error {
	query := `DELETE FROM topics WHERE id = $1`

	result, err := s.db.Exec(query, topicID)
	if err != nil {
		return fmt.Errorf("failed to delete topic: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("topic not found")
	}

	return nil
}

// GetTopicCategories returns all unique categories
func (s *TopicService) GetTopicCategories() ([]string, error) {
	query := `
		SELECT DISTINCT category 
		FROM topics 
		WHERE category IS NOT NULL 
		ORDER BY category
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate categories: %w", err)
	}

	return categories, nil
}