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
	FolderType  string     `json:"folder_type" db:"folder_type"` // "normal", "front_desk", etc.
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`

	// Relationships
	Parent   *Folder   `json:"parent,omitempty"`
	Children []Folder  `json:"children,omitempty"`
	Contents []Content `json:"contents,omitempty"`
	User     *User     `json:"user,omitempty"`
}

// IsFrontDesk returns true if this is the special Front Desk folder
func (f *Folder) IsFrontDesk() bool {
	return f.FolderType == "front_desk"
}

// IsEditable returns true if this folder can be edited
func (f *Folder) IsEditable() bool {
	return f.FolderType == "normal" || f.FolderType == ""
}

// IsDeletable returns true if this folder can be deleted
func (f *Folder) IsDeletable() bool {
	return f.FolderType == "normal" || f.FolderType == ""
}

// Content represents content within a folder
type Content struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Question    string    `json:"question" db:"question"`       // AI-generated question
	Body        string    `json:"body" db:"body"`               // User's answer content
	ContentType string    `json:"content_type" db:"content_type"`
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

// Topic represents a learning subject/domain
type Topic struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Description     string    `json:"description" db:"description"`
	Category        string    `json:"category" db:"category"`
	DifficultyLevel string    `json:"difficulty_level" db:"difficulty_level"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// LearningMethod represents a predefined learning pathway
type LearningMethod struct {
	ID               uuid.UUID `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	Description      string    `json:"description" db:"description"`
	Icon             string    `json:"icon" db:"icon"`
	MethodKey        string    `json:"method_key" db:"method_key"`
	AiPromptTemplate string    `json:"ai_prompt_template" db:"ai_prompt_template"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// UserLearningProgress tracks user progress on topics with specific methods
type UserLearningProgress struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	UserID           uuid.UUID  `json:"user_id" db:"user_id"`
	TopicID          uuid.UUID  `json:"topic_id" db:"topic_id"`
	LearningMethodID uuid.UUID  `json:"learning_method_id" db:"learning_method_id"`
	ProgressPercent  int        `json:"progress_percentage" db:"progress_percentage"`
	Status           string     `json:"status" db:"status"`
	LastActivityAt   time.Time  `json:"last_activity_at" db:"last_activity_at"`
	CompletedAt      *time.Time `json:"completed_at" db:"completed_at"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`

	// Relationships
	User           *User            `json:"user,omitempty"`
	Topic          *Topic           `json:"topic,omitempty"`
	LearningMethod *LearningMethod  `json:"learning_method,omitempty"`
}

// LearningPrompt stores user prompts and AI responses
type LearningPrompt struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	UserID           uuid.UUID  `json:"user_id" db:"user_id"`
	TopicID          uuid.UUID  `json:"topic_id" db:"topic_id"`
	LearningMethodID uuid.UUID  `json:"learning_method_id" db:"learning_method_id"`
	UserPrompt       string     `json:"user_prompt" db:"user_prompt"`
	AiPrompt         string     `json:"ai_prompt" db:"ai_prompt"`
	AiResponse       *string    `json:"ai_response" db:"ai_response"`
	ResponseStatus   string     `json:"response_status" db:"response_status"`
	ErrorMessage     *string    `json:"error_message" db:"error_message"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`

	// Relationships
	User           *User           `json:"user,omitempty"`
	Topic          *Topic          `json:"topic,omitempty"`
	LearningMethod *LearningMethod `json:"learning_method,omitempty"`
}

// LearningSession tracks learning activity
type LearningSession struct {
	ID                   uuid.UUID  `json:"id" db:"id"`
	UserID               uuid.UUID  `json:"user_id" db:"user_id"`
	TopicID              uuid.UUID  `json:"topic_id" db:"topic_id"`
	LearningMethodID     uuid.UUID  `json:"learning_method_id" db:"learning_method_id"`
	SessionDurationMins  int        `json:"session_duration_minutes" db:"session_duration_minutes"`
	ActivitiesCompleted  int        `json:"activities_completed" db:"activities_completed"`
	SessionNotes         *string    `json:"session_notes" db:"session_notes"`
	StartedAt            time.Time  `json:"started_at" db:"started_at"`
	EndedAt              *time.Time `json:"ended_at" db:"ended_at"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`

	// Relationships
	User           *User           `json:"user,omitempty"`
	Topic          *Topic          `json:"topic,omitempty"`
	LearningMethod *LearningMethod `json:"learning_method,omitempty"`
}

// UserPreferences stores user settings and customizations
type UserPreferences struct {
	ID                     uuid.UUID   `json:"id" db:"id"`
	UserID                 uuid.UUID   `json:"user_id" db:"user_id"`
	LearningIntensity      int         `json:"learning_intensity" db:"learning_intensity"`
	SessionDurationMins    int         `json:"session_duration_minutes" db:"session_duration_minutes"`
	DifficultyPreference   int         `json:"difficulty_preference" db:"difficulty_preference"`
	PreferredLearningTimes interface{} `json:"preferred_learning_times" db:"preferred_learning_times"`
	NotificationSettings   interface{} `json:"notification_settings" db:"notification_settings"`
	ThemePreference        string      `json:"theme_preference" db:"theme_preference"`
	CreatedAt              time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time   `json:"updated_at" db:"updated_at"`

	// Relationships
	User *User `json:"user,omitempty"`
}

// LearningPathwayCard represents a pathway card with progress for UI
type LearningPathwayCard struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Icon             string    `json:"icon"`
	MethodKey        string    `json:"method_key"`
	ProgressPercent  int       `json:"progress"`
	Status           string    `json:"status"`
	HasStarted       bool      `json:"has_started"`
	LastActivityAt   *time.Time `json:"last_activity_at"`
}

// DashboardStats represents user dashboard statistics
type DashboardStats struct {
	OverallProgress     int `json:"overall_progress"`
	WeeklyProgress      int `json:"weekly_progress"`
	DailyStreak         int `json:"daily_streak"`
	TotalStudyTime      int `json:"total_study_time"`
	ConceptsMastered    int `json:"concepts_mastered"`
	QuizAccuracy        int `json:"quiz_accuracy"`
	RetentionRate       int `json:"retention_rate"`
}
