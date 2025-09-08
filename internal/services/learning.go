package services

import (
	"database/sql"
	"fmt"
	"mh-go/internal/models"
	"strings"
	"time"

	"github.com/google/uuid"
)

type LearningService struct {
	db                   *sql.DB
	topicService        *TopicService
	methodService       *LearningMethodService
	progressService     *ProgressService
}

func NewLearningService(db *sql.DB, topicService *TopicService, methodService *LearningMethodService, progressService *ProgressService) *LearningService {
	return &LearningService{
		db:              db,
		topicService:    topicService,
		methodService:   methodService,
		progressService: progressService,
	}
}

// CreateLearningPrompt creates a new learning prompt with AI generation
func (s *LearningService) CreateLearningPrompt(userID, topicID, methodID uuid.UUID, userPrompt string) (*models.LearningPrompt, error) {
	// Get the topic to use its name in the AI prompt
	topic, err := s.topicService.GetTopicByID(topicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}

	// Generate the AI prompt using the method template
	aiPrompt, err := s.methodService.GenerateAIPrompt(methodID, topic.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI prompt: %w", err)
	}

	// Append the user's specific prompt to the AI prompt
	fullAIPrompt := aiPrompt + "\n\n" + "User's specific learning goal: " + userPrompt

	prompt := &models.LearningPrompt{
		ID:               uuid.New(),
		UserID:           userID,
		TopicID:          topicID,
		LearningMethodID: methodID,
		UserPrompt:       userPrompt,
		AiPrompt:         fullAIPrompt,
		ResponseStatus:   "pending",
	}

	query := `
		INSERT INTO learning_prompts 
		(id, user_id, topic_id, learning_method_id, user_prompt, ai_prompt, response_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`

	err = s.db.QueryRow(
		query, prompt.ID, prompt.UserID, prompt.TopicID,
		prompt.LearningMethodID, prompt.UserPrompt, prompt.AiPrompt,
		prompt.ResponseStatus,
	).Scan(&prompt.CreatedAt, &prompt.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create learning prompt: %w", err)
	}

	// Create or update progress to mark as started
	_, err = s.progressService.CreateOrUpdateProgress(userID, topicID, methodID, 10, "in_progress")
	if err != nil {
		// Log error but don't fail the prompt creation
		fmt.Printf("Warning: failed to update progress: %v\n", err)
	}

	return prompt, nil
}

// UpdatePromptResponse updates the AI response for a prompt
func (s *LearningService) UpdatePromptResponse(promptID uuid.UUID, response string, status string) error {
	query := `
		UPDATE learning_prompts
		SET ai_response = $2, response_status = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := s.db.Exec(query, promptID, response, status)
	if err != nil {
		return fmt.Errorf("failed to update prompt response: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("prompt not found")
	}

	return nil
}

// UpdatePromptError updates the error message for a failed prompt
func (s *LearningService) UpdatePromptError(promptID uuid.UUID, errorMessage string) error {
	query := `
		UPDATE learning_prompts
		SET error_message = $2, response_status = 'error', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := s.db.Exec(query, promptID, errorMessage)
	if err != nil {
		return fmt.Errorf("failed to update prompt error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("prompt not found")
	}

	return nil
}

// GetUserPrompts returns all prompts for a user
func (s *LearningService) GetUserPrompts(userID uuid.UUID, limit int) ([]models.LearningPrompt, error) {
	query := `
		SELECT id, user_id, topic_id, learning_method_id, user_prompt, ai_prompt, 
		       ai_response, response_status, error_message, created_at, updated_at
		FROM learning_prompts
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := s.db.Query(query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query user prompts: %w", err)
	}
	defer rows.Close()

	var prompts []models.LearningPrompt
	for rows.Next() {
		var prompt models.LearningPrompt
		var aiResponse, errorMessage sql.NullString

		err := rows.Scan(
			&prompt.ID, &prompt.UserID, &prompt.TopicID, &prompt.LearningMethodID,
			&prompt.UserPrompt, &prompt.AiPrompt, &aiResponse, &prompt.ResponseStatus,
			&errorMessage, &prompt.CreatedAt, &prompt.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan prompt: %w", err)
		}

		if aiResponse.Valid {
			prompt.AiResponse = &aiResponse.String
		}
		if errorMessage.Valid {
			prompt.ErrorMessage = &errorMessage.String
		}

		prompts = append(prompts, prompt)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate prompts: %w", err)
	}

	return prompts, nil
}

// GetPromptByID returns a prompt by ID
func (s *LearningService) GetPromptByID(promptID uuid.UUID) (*models.LearningPrompt, error) {
	var prompt models.LearningPrompt
	var aiResponse, errorMessage sql.NullString

	query := `
		SELECT id, user_id, topic_id, learning_method_id, user_prompt, ai_prompt, 
		       ai_response, response_status, error_message, created_at, updated_at
		FROM learning_prompts
		WHERE id = $1
	`

	err := s.db.QueryRow(query, promptID).Scan(
		&prompt.ID, &prompt.UserID, &prompt.TopicID, &prompt.LearningMethodID,
		&prompt.UserPrompt, &prompt.AiPrompt, &aiResponse, &prompt.ResponseStatus,
		&errorMessage, &prompt.CreatedAt, &prompt.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("prompt not found")
		}
		return nil, fmt.Errorf("failed to get prompt: %w", err)
	}

	if aiResponse.Valid {
		prompt.AiResponse = &aiResponse.String
	}
	if errorMessage.Valid {
		prompt.ErrorMessage = &errorMessage.String
	}

	return &prompt, nil
}

// StartLearningSession creates a new learning session
func (s *LearningService) StartLearningSession(userID, topicID, methodID uuid.UUID) (*models.LearningSession, error) {
	session := &models.LearningSession{
		ID:               uuid.New(),
		UserID:           userID,
		TopicID:          topicID,
		LearningMethodID: methodID,
		StartedAt:        time.Now(),
	}

	query := `
		INSERT INTO learning_sessions 
		(id, user_id, topic_id, learning_method_id, started_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at
	`

	err := s.db.QueryRow(
		query, session.ID, session.UserID, session.TopicID,
		session.LearningMethodID, session.StartedAt,
	).Scan(&session.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create learning session: %w", err)
	}

	return session, nil
}

// EndLearningSession completes a learning session
func (s *LearningService) EndLearningSession(sessionID uuid.UUID, activitiesCompleted int, notes string) error {
	now := time.Now()

	query := `
		UPDATE learning_sessions
		SET ended_at = $2, 
		    session_duration_minutes = EXTRACT(EPOCH FROM ($2 - started_at)) / 60,
		    activities_completed = $3,
		    session_notes = $4
		WHERE id = $1
	`

	result, err := s.db.Exec(query, sessionID, now, activitiesCompleted, notes)
	if err != nil {
		return fmt.Errorf("failed to end learning session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

// GetUserSessions returns recent learning sessions for a user
func (s *LearningService) GetUserSessions(userID uuid.UUID, limit int) ([]models.LearningSession, error) {
	query := `
		SELECT id, user_id, topic_id, learning_method_id, session_duration_minutes,
		       activities_completed, session_notes, started_at, ended_at, created_at
		FROM learning_sessions
		WHERE user_id = $1
		ORDER BY started_at DESC
		LIMIT $2
	`

	rows, err := s.db.Query(query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query user sessions: %w", err)
	}
	defer rows.Close()

	var sessions []models.LearningSession
	for rows.Next() {
		var session models.LearningSession
		var endedAt sql.NullTime
		var sessionNotes sql.NullString

		err := rows.Scan(
			&session.ID, &session.UserID, &session.TopicID, &session.LearningMethodID,
			&session.SessionDurationMins, &session.ActivitiesCompleted, &sessionNotes,
			&session.StartedAt, &endedAt, &session.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}

		if endedAt.Valid {
			session.EndedAt = &endedAt.Time
		}
		if sessionNotes.Valid {
			session.SessionNotes = &sessionNotes.String
		}

		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate sessions: %w", err)
	}

	return sessions, nil
}

// SearchTopics searches for topics by name or category
func (s *LearningService) SearchTopics(query string, limit int) ([]models.Topic, error) {
	searchQuery := `
		SELECT id, name, description, category, difficulty_level, created_at, updated_at
		FROM topics
		WHERE name ILIKE '%' || $1 || '%' 
		   OR description ILIKE '%' || $1 || '%'
		   OR category ILIKE '%' || $1 || '%'
		ORDER BY 
			CASE WHEN name ILIKE $1 || '%' THEN 1
			     WHEN name ILIKE '%' || $1 || '%' THEN 2
			     ELSE 3 END,
			name
		LIMIT $2
	`

	rows, err := s.db.Query(searchQuery, strings.ToLower(query), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search topics: %w", err)
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

// GetOrCreateUserPreferences gets user preferences or creates default ones
func (s *LearningService) GetOrCreateUserPreferences(userID uuid.UUID) (*models.UserPreferences, error) {
	// Try to get existing preferences
	var prefs models.UserPreferences
	query := `
		SELECT id, user_id, learning_intensity, session_duration_minutes, difficulty_preference,
		       preferred_learning_times, notification_settings, theme_preference, created_at, updated_at
		FROM user_preferences
		WHERE user_id = $1
	`

	err := s.db.QueryRow(query, userID).Scan(
		&prefs.ID, &prefs.UserID, &prefs.LearningIntensity, &prefs.SessionDurationMins,
		&prefs.DifficultyPreference, &prefs.PreferredLearningTimes, &prefs.NotificationSettings,
		&prefs.ThemePreference, &prefs.CreatedAt, &prefs.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Create default preferences
			return s.createDefaultUserPreferences(userID)
		}
		return nil, fmt.Errorf("failed to get user preferences: %w", err)
	}

	return &prefs, nil
}

// createDefaultUserPreferences creates default preferences for a new user
func (s *LearningService) createDefaultUserPreferences(userID uuid.UUID) (*models.UserPreferences, error) {
	prefs := &models.UserPreferences{
		ID:                   uuid.New(),
		UserID:               userID,
		LearningIntensity:    7,
		SessionDurationMins:  45,
		DifficultyPreference: 6,
		ThemePreference:      "system",
	}

	query := `
		INSERT INTO user_preferences 
		(id, user_id, learning_intensity, session_duration_minutes, difficulty_preference, theme_preference)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`

	err := s.db.QueryRow(
		query, prefs.ID, prefs.UserID, prefs.LearningIntensity,
		prefs.SessionDurationMins, prefs.DifficultyPreference, prefs.ThemePreference,
	).Scan(&prefs.CreatedAt, &prefs.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create default preferences: %w", err)
	}

	return prefs, nil
}

// UpdateUserPreferences updates user preferences
func (s *LearningService) UpdateUserPreferences(prefs *models.UserPreferences) error {
	query := `
		UPDATE user_preferences
		SET learning_intensity = $2, session_duration_minutes = $3, difficulty_preference = $4,
		    preferred_learning_times = $5, notification_settings = $6, theme_preference = $7,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`

	err := s.db.QueryRow(
		query, prefs.ID, prefs.LearningIntensity, prefs.SessionDurationMins,
		prefs.DifficultyPreference, prefs.PreferredLearningTimes, prefs.NotificationSettings,
		prefs.ThemePreference,
	).Scan(&prefs.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("preferences not found")
		}
		return fmt.Errorf("failed to update preferences: %w", err)
	}

	return nil
}