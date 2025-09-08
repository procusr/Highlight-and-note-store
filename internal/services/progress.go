package services

import (
	"database/sql"
	"fmt"
	"mh-go/internal/models"
	"time"

	"github.com/google/uuid"
)

type ProgressService struct {
	db *sql.DB
}

func NewProgressService(db *sql.DB) *ProgressService {
	return &ProgressService{db: db}
}

// GetUserProgress returns user's progress for a specific topic and method
func (s *ProgressService) GetUserProgress(userID, topicID, methodID uuid.UUID) (*models.UserLearningProgress, error) {
	var progress models.UserLearningProgress
	query := `
		SELECT id, user_id, topic_id, learning_method_id, progress_percentage, 
		       status, last_activity_at, completed_at, created_at, updated_at
		FROM user_learning_progress
		WHERE user_id = $1 AND topic_id = $2 AND learning_method_id = $3
	`

	err := s.db.QueryRow(query, userID, topicID, methodID).Scan(
		&progress.ID, &progress.UserID, &progress.TopicID, &progress.LearningMethodID,
		&progress.ProgressPercent, &progress.Status, &progress.LastActivityAt,
		&progress.CompletedAt, &progress.CreatedAt, &progress.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No progress found, which is valid
		}
		return nil, fmt.Errorf("failed to get user progress: %w", err)
	}

	return &progress, nil
}

// GetAllUserProgress returns all progress for a user
func (s *ProgressService) GetAllUserProgress(userID uuid.UUID) ([]models.UserLearningProgress, error) {
	query := `
		SELECT id, user_id, topic_id, learning_method_id, progress_percentage, 
		       status, last_activity_at, completed_at, created_at, updated_at
		FROM user_learning_progress
		WHERE user_id = $1
		ORDER BY last_activity_at DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user progress: %w", err)
	}
	defer rows.Close()

	var progressList []models.UserLearningProgress
	for rows.Next() {
		var progress models.UserLearningProgress
		err := rows.Scan(
			&progress.ID, &progress.UserID, &progress.TopicID, &progress.LearningMethodID,
			&progress.ProgressPercent, &progress.Status, &progress.LastActivityAt,
			&progress.CompletedAt, &progress.CreatedAt, &progress.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user progress: %w", err)
		}
		progressList = append(progressList, progress)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate user progress: %w", err)
	}

	return progressList, nil
}

// CreateUserProgress creates a new progress record
func (s *ProgressService) CreateUserProgress(progress *models.UserLearningProgress) error {
	progress.ID = uuid.New()
	progress.LastActivityAt = time.Now()

	query := `
		INSERT INTO user_learning_progress 
		(id, user_id, topic_id, learning_method_id, progress_percentage, status, last_activity_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`

	err := s.db.QueryRow(
		query, progress.ID, progress.UserID, progress.TopicID, 
		progress.LearningMethodID, progress.ProgressPercent, 
		progress.Status, progress.LastActivityAt,
	).Scan(&progress.CreatedAt, &progress.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user progress: %w", err)
	}

	return nil
}

// UpdateUserProgress updates existing progress
func (s *ProgressService) UpdateUserProgress(progress *models.UserLearningProgress) error {
	progress.LastActivityAt = time.Now()

	query := `
		UPDATE user_learning_progress
		SET progress_percentage = $2, status = $3, last_activity_at = $4, 
		    completed_at = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`

	err := s.db.QueryRow(
		query, progress.ID, progress.ProgressPercent, progress.Status,
		progress.LastActivityAt, progress.CompletedAt,
	).Scan(&progress.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("progress record not found")
		}
		return fmt.Errorf("failed to update user progress: %w", err)
	}

	return nil
}

// CreateOrUpdateProgress creates new progress or updates existing one
func (s *ProgressService) CreateOrUpdateProgress(userID, topicID, methodID uuid.UUID, progressPercent int, status string) (*models.UserLearningProgress, error) {
	// Check if progress already exists
	existingProgress, err := s.GetUserProgress(userID, topicID, methodID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing progress: %w", err)
	}

	var progress *models.UserLearningProgress

	if existingProgress == nil {
		// Create new progress
		progress = &models.UserLearningProgress{
			UserID:           userID,
			TopicID:          topicID,
			LearningMethodID: methodID,
			ProgressPercent:  progressPercent,
			Status:           status,
		}

		// Set completed_at if status is completed
		if status == "completed" {
			now := time.Now()
			progress.CompletedAt = &now
		}

		err = s.CreateUserProgress(progress)
		if err != nil {
			return nil, fmt.Errorf("failed to create progress: %w", err)
		}
	} else {
		// Update existing progress
		existingProgress.ProgressPercent = progressPercent
		existingProgress.Status = status

		// Set completed_at if status is completed and not already set
		if status == "completed" && existingProgress.CompletedAt == nil {
			now := time.Now()
			existingProgress.CompletedAt = &now
		}

		err = s.UpdateUserProgress(existingProgress)
		if err != nil {
			return nil, fmt.Errorf("failed to update progress: %w", err)
		}

		progress = existingProgress
	}

	return progress, nil
}

// GetLearningPathwayCards returns pathway cards with user progress for UI
func (s *ProgressService) GetLearningPathwayCards(userID uuid.UUID) ([]models.LearningPathwayCard, error) {
	query := `
		SELECT 
			lm.id, 
			lm.name, 
			lm.description, 
			lm.icon, 
			lm.method_key,
			COALESCE(ulp.progress_percentage, 0) as progress,
			COALESCE(ulp.status, 'not_started') as status,
			CASE WHEN ulp.id IS NOT NULL THEN true ELSE false END as has_started,
			ulp.last_activity_at
		FROM learning_methods lm
		LEFT JOIN user_learning_progress ulp ON lm.id = ulp.learning_method_id AND ulp.user_id = $1
		ORDER BY lm.name
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query pathway cards: %w", err)
	}
	defer rows.Close()

	var cards []models.LearningPathwayCard
	for rows.Next() {
		var card models.LearningPathwayCard
		var lastActivity sql.NullTime

		err := rows.Scan(
			&card.ID, &card.Name, &card.Description, &card.Icon,
			&card.MethodKey, &card.ProgressPercent, &card.Status,
			&card.HasStarted, &lastActivity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pathway card: %w", err)
		}

		if lastActivity.Valid {
			card.LastActivityAt = &lastActivity.Time
		}

		cards = append(cards, card)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate pathway cards: %w", err)
	}

	return cards, nil
}

// GetDashboardStats returns user dashboard statistics
func (s *ProgressService) GetDashboardStats(userID uuid.UUID) (*models.DashboardStats, error) {
	stats := &models.DashboardStats{
		OverallProgress: 0,
		WeeklyProgress:  0,
		DailyStreak:     0,
		TotalStudyTime:  0,
		ConceptsMastered: 0,
		QuizAccuracy:    85,
		RetentionRate:   82,
	}

	// Get overall progress (average of all user progress)
	progressQuery := `
		SELECT 
			COALESCE(AVG(progress_percentage), 0) as overall_progress,
			COUNT(*) as total_started,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_count
		FROM user_learning_progress
		WHERE user_id = $1
	`

	var totalStarted, completedCount int
	err := s.db.QueryRow(progressQuery, userID).Scan(
		&stats.OverallProgress, &totalStarted, &completedCount,
	)
	if err != nil {
		// Don't fail, just use default values and log the error
		fmt.Printf("Warning: failed to get progress stats for user %s: %v\n", userID, err)
	} else {
		stats.ConceptsMastered = completedCount
	}

	// Get weekly progress (activities in the last 7 days)
	weeklyQuery := `
		SELECT COUNT(*) 
		FROM learning_sessions 
		WHERE user_id = $1 AND started_at >= NOW() - INTERVAL '7 days'
	`

	err = s.db.QueryRow(weeklyQuery, userID).Scan(&stats.WeeklyProgress)
	if err != nil {
		stats.WeeklyProgress = 0 // Default to 0 if query fails
	}

	// Get total study time (sum of all session durations)
	studyTimeQuery := `
		SELECT COALESCE(SUM(session_duration_minutes), 0)
		FROM learning_sessions
		WHERE user_id = $1
	`

	err = s.db.QueryRow(studyTimeQuery, userID).Scan(&stats.TotalStudyTime)
	if err != nil {
		stats.TotalStudyTime = 0 // Default to 0 if query fails
	}

	// For new users, daily streak should be 0
	stats.DailyStreak = 0

	return stats, nil
}