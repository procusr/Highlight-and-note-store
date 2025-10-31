package services

import (
	"context"
	"fmt"
	"mh-go/internal/models"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Enhanced AI-Powered Learning Methods

// CreateEnhancedLearningSession creates a new session with initial AI content generation
func (s *LearningService) CreateEnhancedLearningSession(ctx context.Context, req models.CreateSessionRequest, userID uuid.UUID) (*models.SessionResponse, error) {
	// Create a topic for the user's learning session
	topic := &models.Topic{
		Name:            req.TopicName,
		Description:     fmt.Sprintf("Learning topic: %s", req.TopicName),
		Category:        "user-generated",
		DifficultyLevel: "intermediate",
	}

	err := s.topicService.CreateTopic(topic)
	if err != nil {
		return nil, fmt.Errorf("failed to create topic: %w", err)
	}

	// Get learning method
	method, err := s.methodService.GetLearningMethodByID(req.LearningMethodID)
	if err != nil {
		return nil, fmt.Errorf("failed to get learning method: %w", err)
	}

	// Create session record in DB
	session := &models.LearningSession{
		ID:               uuid.New(),
		UserID:           userID,
		TopicID:          topic.ID,
		LearningMethodID: req.LearningMethodID,
		StartedAt:        time.Now(),
		Status:           string(models.SessionStatusActive),
		CurrentStep:      1,
	}

	query := `
		INSERT INTO learning_sessions
		(id, user_id, topic_id, learning_method_id, started_at, status, current_step)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at
	`

	err = s.db.QueryRow(
		query, session.ID, session.UserID, session.TopicID,
		session.LearningMethodID, session.StartedAt, session.Status, session.CurrentStep,
	).Scan(&session.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create learning session: %w", err)
	}

	// Build system prompt with variable replacement
	systemPrompt := s.buildSystemPrompt(method.AiPromptTemplate, req.TopicName, req.UserGoal)

	// Generate initial AI content synchronously
	aiRequest := models.AIRequest{
		SystemPrompt:    systemPrompt,
		CurrentInput:    "Begin the learning session with an overview and first lesson.",
		InteractionType: models.InteractionInitialContent,
	}

	initialContent, err := s.openaiService.GenerateLearningContent(ctx, aiRequest)
	if err != nil {
		// Rollback session creation if AI generation fails
		s.db.Exec("DELETE FROM learning_sessions WHERE id = $1", session.ID)
		return nil, fmt.Errorf("failed to generate initial content: %w", err)
	}

	// Store initial interaction
	stepNumber := 1
	interactionType := string(models.InteractionInitialContent)
	prompt := &models.LearningPrompt{
		ID:               uuid.New(),
		UserID:           userID,
		TopicID:          topic.ID,
		LearningMethodID: req.LearningMethodID,
		LearningGoal:     req.UserGoal,
		AiPrompt:         systemPrompt,
		AiResponse:       &initialContent,
		ResponseStatus:   "success",
		SessionID:        &session.ID,
		StepNumber:       &stepNumber,
		InteractionType:  &interactionType,
	}

	promptQuery := `
		INSERT INTO learning_prompts
		(id, user_id, topic_id, learning_method_id, learning_goal, ai_prompt, ai_response, 
		 response_status, session_id, step_number, interaction_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING created_at, updated_at
	`

	err = s.db.QueryRow(
		promptQuery, prompt.ID, prompt.UserID, prompt.TopicID, prompt.LearningMethodID,
		prompt.LearningGoal, prompt.AiPrompt, prompt.AiResponse, prompt.ResponseStatus,
		prompt.SessionID, prompt.StepNumber, prompt.InteractionType,
	).Scan(&prompt.CreatedAt, &prompt.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to store initial interaction: %w", err)
	}

	return &models.SessionResponse{
		SessionID:      session.ID,
		InitialContent: initialContent,
		LearningPath:   []string{}, // TODO: Generate learning path based on topic
	}, nil
}

// ProcessSessionInteraction handles user interactions within a learning session
func (s *LearningService) ProcessSessionInteraction(ctx context.Context, sessionID uuid.UUID, req models.InteractionRequest) (*models.AIResponse, error) {
	// Load session and verify it's active
	session, err := s.GetLearningSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if session.Status != string(models.SessionStatusActive) {
		return nil, fmt.Errorf("session is not active")
	}

	// Get learning method for system prompt
	method, err := s.methodService.GetLearningMethodByID(session.LearningMethodID)
	if err != nil {
		return nil, fmt.Errorf("failed to get learning method: %w", err)
	}

	// Get recent interactions (sliding window)
	recentInteractions, err := s.getRecentInteractions(sessionID, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent interactions: %w", err)
	}

	// Build system prompt
	systemPrompt := s.buildSystemPrompt(method.AiPromptTemplate, session.Topic.Name, session.Goal)

	// Build AI request
	aiRequest := models.AIRequest{
		SystemPrompt:       systemPrompt,
		SessionSummary:     s.getSessionSummary(session),
		RecentInteractions: recentInteractions,
		CurrentInput:       req.Message,
		InteractionType:    req.InteractionType,
	}

	// Generate AI response
	response, err := s.openaiService.GenerateLearningContent(ctx, aiRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI response: %w", err)
	}

	// Save interaction to DB (auto-save)
	err = s.saveInteraction(sessionID, req.Message, response, req.InteractionType, session.CurrentStep)
	if err != nil {
		return nil, fmt.Errorf("failed to save interaction: %w", err)
	}

	// Update session activity
	s.updateSessionActivity(sessionID)

	return &models.AIResponse{
		Content:           response,
		ContextualActions: s.determineContextualActions(response, session),
		SuggestedActions:  []string{}, // TODO: Implement suggested actions
	}, nil
}

// CompleteConceptAndSummarize marks a concept as complete and generates a summary
func (s *LearningService) CompleteConceptAndSummarize(ctx context.Context, sessionID uuid.UUID, conceptName string) error {
	// Get session
	session, err := s.GetLearningSessionByID(sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	// Get interactions for this concept
	interactions, err := s.getConceptInteractions(sessionID, conceptName)
	if err != nil {
		return fmt.Errorf("failed to get concept interactions: %w", err)
	}

	// Generate AI summary
	summary, err := s.openaiService.GenerateSessionSummary(ctx, interactions, conceptName)
	if err != nil {
		return fmt.Errorf("failed to generate summary: %w", err)
	}

	// Update session with completed concept and summary
	updatedSummary := s.appendToSessionSummary(session.SessionSummary, conceptName, summary)
	updatedConcepts := append(session.ConceptsCompleted, conceptName)

	query := `
		UPDATE learning_sessions
		SET session_summary = $2,
		    concepts_completed = $3,
		    current_step = current_step + 1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err = s.db.Exec(query, sessionID, updatedSummary, updatedConcepts)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	return nil
}

// PauseSession pauses an active learning session
func (s *LearningService) PauseSession(sessionID uuid.UUID) error {
	query := `
		UPDATE learning_sessions
		SET status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND status = $3
	`

	result, err := s.db.Exec(query, sessionID, string(models.SessionStatusPaused), string(models.SessionStatusActive))
	if err != nil {
		return fmt.Errorf("failed to pause session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found or not active")
	}

	return nil
}

// ResumeSession resumes a paused learning session
func (s *LearningService) ResumeSession(sessionID uuid.UUID) error {
	query := `
		UPDATE learning_sessions
		SET status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND status = $3
	`

	result, err := s.db.Exec(query, sessionID, string(models.SessionStatusActive), string(models.SessionStatusPaused))
	if err != nil {
		return fmt.Errorf("failed to resume session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found or not paused")
	}

	return nil
}

// Helper methods

func (s *LearningService) buildSystemPrompt(template, topicName, userGoal string) string {
	prompt := strings.ReplaceAll(template, "{TOPIC}", topicName)
	prompt = strings.ReplaceAll(prompt, "{USER_GOAL}", userGoal)
	return prompt
}

func (s *LearningService) getRecentInteractions(sessionID uuid.UUID, limit int) ([]models.LearningInteraction, error) {
	query := `
		SELECT ai_prompt, ai_response, interaction_type
		FROM learning_prompts
		WHERE session_id = $1 AND ai_response IS NOT NULL
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := s.db.Query(query, sessionID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent interactions: %w", err)
	}
	defer rows.Close()

	var interactions []models.LearningInteraction
	for rows.Next() {
		var aiPrompt, aiResponse, interactionType string
		err := rows.Scan(&aiPrompt, &aiResponse, &interactionType)
		if err != nil {
			return nil, fmt.Errorf("failed to scan interaction: %w", err)
		}

		// Add user message if it's not initial content
		if interactionType != string(models.InteractionInitialContent) {
			interactions = append(interactions, models.LearningInteraction{
				Role:    "user",
				Content: "Previous user input", // TODO: Store actual user input
			})
		}

		// Add AI response
		interactions = append(interactions, models.LearningInteraction{
			Role:    "assistant",
			Content: aiResponse,
		})
	}

	// Reverse to get chronological order
	for i, j := 0, len(interactions)-1; i < j; i, j = i+1, j-1 {
		interactions[i], interactions[j] = interactions[j], interactions[i]
	}

	return interactions, nil
}

func (s *LearningService) getSessionSummary(session *models.LearningSession) string {
	if session.SessionSummary != nil {
		return *session.SessionSummary
	}
	return ""
}

func (s *LearningService) saveInteraction(sessionID uuid.UUID, userInput, aiResponse string, interactionType models.InteractionType, stepNumber int) error {
	query := `
		INSERT INTO learning_prompts
		(id, user_id, topic_id, learning_method_id, learning_goal, ai_prompt, ai_response,
		 response_status, session_id, step_number, interaction_type)
		SELECT $1, user_id, topic_id, learning_method_id, 'Ongoing session', $2, $3, 'success', $4, $5, $6
		FROM learning_sessions WHERE id = $4
	`

	_, err := s.db.Exec(query, uuid.New(), userInput, aiResponse, sessionID, stepNumber, string(interactionType))
	return err
}

func (s *LearningService) updateSessionActivity(sessionID uuid.UUID) error {
	query := `UPDATE learning_sessions SET updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := s.db.Exec(query, sessionID)
	return err
}

func (s *LearningService) determineContextualActions(response string, session *models.LearningSession) []models.ContextualAction {
	// TODO: Implement contextual action logic based on AI response content
	return []models.ContextualAction{
		{
			Type:        "continue",
			Label:       "Continue",
			Action:      "/api/learning/session/" + session.ID.String() + "/next-step",
			Style:       "btn-primary",
			Enabled:     true,
			Description: "Move to the next step",
		},
		{
			Type:        "pause",
			Label:       "Pause Session",
			Action:      "/api/learning/session/" + session.ID.String() + "/pause",
			Style:       "btn-secondary",
			Enabled:     true,
			Description: "Save and pause this session",
		},
	}
}

func (s *LearningService) getConceptInteractions(sessionID uuid.UUID, conceptName string) ([]models.LearningInteraction, error) {
	// For now, return all interactions for the session
	// TODO: Implement concept-specific filtering
	return s.getRecentInteractions(sessionID, 100)
}

func (s *LearningService) appendToSessionSummary(existingSummary *string, conceptName, summary string) string {
	base := ""
	if existingSummary != nil {
		base = *existingSummary
	}

	if base != "" {
		base += "\n\n"
	}

	base += fmt.Sprintf("### Concept: %s\n%s", conceptName, summary)
	return base
}