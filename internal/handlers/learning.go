package handlers

import (
	"encoding/json"
	"fmt"
	"mh-go/internal/middleware"
	"mh-go/internal/services"
	"mh-go/templates"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type LearningHandler struct {
	topicService      *services.TopicService
	methodService     *services.LearningMethodService
	progressService   *services.ProgressService
	learningService   *services.LearningService
	userService       *services.UserService
}

func NewLearningHandler(
	topicService *services.TopicService,
	methodService *services.LearningMethodService,
	progressService *services.ProgressService,
	learningService *services.LearningService,
	userService *services.UserService,
) *LearningHandler {
	return &LearningHandler{
		topicService:    topicService,
		methodService:   methodService,
		progressService: progressService,
		learningService: learningService,
		userService:     userService,
	}
}

// LearningPlatformPage renders the main learning platform page
func (h *LearningHandler) LearningPlatformPage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	// Get user model
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Get dashboard stats
	stats, err := h.progressService.GetDashboardStats(userID)
	if err != nil {
		http.Error(w, "Failed to get dashboard stats", http.StatusInternalServerError)
		return
	}

	// Get learning pathway cards
	pathwayCards, err := h.progressService.GetLearningPathwayCards(userID)
	if err != nil {
		http.Error(w, "Failed to get learning pathways", http.StatusInternalServerError)
		return
	}

	// Get topics
	topics, err := h.topicService.GetAllTopics()
	if err != nil {
		http.Error(w, "Failed to get topics", http.StatusInternalServerError)
		return
	}

	// Get active learning sessions
	sessions, err := h.learningService.GetActiveLearningSessions(userID)
	if err != nil {
		http.Error(w, "Failed to get learning sessions", http.StatusInternalServerError)
		return
	}

	// Get user preferences
	preferences, err := h.learningService.GetOrCreateUserPreferences(userID)
	if err != nil {
		http.Error(w, "Failed to get user preferences", http.StatusInternalServerError)
		return
	}

	component := templates.LearningPlatformPage(user, stats, pathwayCards, topics, sessions, preferences)
	component.Render(r.Context(), w)
}

// GetLearningPathways returns learning pathway cards for HTMX
func (h *LearningHandler) GetLearningPathways(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	pathwayCards, err := h.progressService.GetLearningPathwayCards(userID)
	if err != nil {
		http.Error(w, "Failed to get learning pathways", http.StatusInternalServerError)
		return
	}

	component := templates.LearningPathwaysGrid(pathwayCards)
	component.Render(r.Context(), w)
}

// GetTopics returns available topics as JSON
func (h *LearningHandler) GetTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := h.topicService.GetAllTopics()
	if err != nil {
		http.Error(w, "Failed to get topics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topics)
}

// SearchTopics searches for topics
func (h *LearningHandler) SearchTopics(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	topics, err := h.learningService.SearchTopics(query, limit)
	if err != nil {
		http.Error(w, "Failed to search topics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topics)
}

// CreateLearningPrompt handles prompt creation from learning method cards
func (h *LearningHandler) CreateLearningPrompt(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req struct {
		TopicID    string `json:"topic_id"`
		MethodID   string `json:"method_id"`
		UserPrompt string `json:"user_prompt"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate UUIDs
	topicID, err := uuid.Parse(req.TopicID)
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	methodID, err := uuid.Parse(req.MethodID)
	if err != nil {
		http.Error(w, "Invalid method ID", http.StatusBadRequest)
		return
	}

	if req.UserPrompt == "" {
		http.Error(w, "User prompt is required", http.StatusBadRequest)
		return
	}

	// Create the learning prompt
	prompt, err := h.learningService.CreateLearningPrompt(userID, topicID, methodID, req.UserPrompt)
	if err != nil {
		http.Error(w, "Failed to create learning prompt", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"prompt_id": prompt.ID,
		"ai_prompt": prompt.AiPrompt,
		"message":   "Learning prompt created successfully",
	})
}

// GetDashboardStats returns dashboard statistics
func (h *LearningHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	stats, err := h.progressService.GetDashboardStats(userID)
	if err != nil {
		http.Error(w, "Failed to get dashboard stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// UpdateProgress updates user progress for a topic/method combination
func (h *LearningHandler) UpdateProgress(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req struct {
		TopicID         string `json:"topic_id"`
		MethodID        string `json:"method_id"`
		ProgressPercent int    `json:"progress_percentage"`
		Status          string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate UUIDs
	topicID, err := uuid.Parse(req.TopicID)
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	methodID, err := uuid.Parse(req.MethodID)
	if err != nil {
		http.Error(w, "Invalid method ID", http.StatusBadRequest)
		return
	}

	// Validate progress percentage
	if req.ProgressPercent < 0 || req.ProgressPercent > 100 {
		http.Error(w, "Progress percentage must be between 0 and 100", http.StatusBadRequest)
		return
	}

	// Update progress
	progress, err := h.progressService.CreateOrUpdateProgress(userID, topicID, methodID, req.ProgressPercent, req.Status)
	if err != nil {
		http.Error(w, "Failed to update progress", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"progress": progress,
		"message":  "Progress updated successfully",
	})
}

// StartLearningSession starts a new learning session
func (h *LearningHandler) StartLearningSession(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req struct {
		TopicID  string `json:"topic_id"`
		MethodID string `json:"method_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate UUIDs
	topicID, err := uuid.Parse(req.TopicID)
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	methodID, err := uuid.Parse(req.MethodID)
	if err != nil {
		http.Error(w, "Invalid method ID", http.StatusBadRequest)
		return
	}

	// Start learning session
	session, err := h.learningService.StartLearningSession(userID, topicID, methodID)
	if err != nil {
		http.Error(w, "Failed to start learning session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"session_id": session.ID,
		"message":    "Learning session started",
	})
}

// EndLearningSession ends a learning session
func (h *LearningHandler) EndLearningSession(w http.ResponseWriter, r *http.Request) {
	sessionID, err := uuid.Parse(chi.URLParam(r, "sessionID"))
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	var req struct {
		ActivitiesCompleted int    `json:"activities_completed"`
		Notes               string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.learningService.EndLearningSession(sessionID, req.ActivitiesCompleted, req.Notes)
	if err != nil {
		http.Error(w, "Failed to end learning session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Learning session ended successfully",
	})
}

// GetUserPrompts returns recent prompts for a user
func (h *LearningHandler) GetUserPrompts(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	prompts, err := h.learningService.GetUserPrompts(userID, limit)
	if err != nil {
		http.Error(w, "Failed to get user prompts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prompts)
}

// GetUserPreferences returns user preferences
func (h *LearningHandler) GetUserPreferences(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	preferences, err := h.learningService.GetOrCreateUserPreferences(userID)
	if err != nil {
		http.Error(w, "Failed to get user preferences", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(preferences)
}

// GetActiveLearningSessions returns active learning sessions for the current user
func (h *LearningHandler) GetActiveLearningSessions(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	sessions, err := h.learningService.GetActiveLearningSessions(userID)
	if err != nil {
		http.Error(w, "Failed to get active learning sessions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"sessions":  sessions,
		"message":   "Active learning sessions retrieved successfully",
	})
}

// UpdateUserPreferences updates user preferences
func (h *LearningHandler) UpdateUserPreferences(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req struct {
		LearningIntensity    int    `json:"learning_intensity"`
		SessionDurationMins  int    `json:"session_duration_minutes"`
		DifficultyPreference int    `json:"difficulty_preference"`
		ThemePreference      string `json:"theme_preference"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Get existing preferences
	preferences, err := h.learningService.GetOrCreateUserPreferences(userID)
	if err != nil {
		http.Error(w, "Failed to get user preferences", http.StatusInternalServerError)
		return
	}

	// Update values
	preferences.LearningIntensity = req.LearningIntensity
	preferences.SessionDurationMins = req.SessionDurationMins
	preferences.DifficultyPreference = req.DifficultyPreference
	preferences.ThemePreference = req.ThemePreference

	// Validate values
	if preferences.LearningIntensity < 1 || preferences.LearningIntensity > 10 {
		http.Error(w, "Learning intensity must be between 1 and 10", http.StatusBadRequest)
		return
	}

	if preferences.SessionDurationMins < 15 || preferences.SessionDurationMins > 90 {
		http.Error(w, "Session duration must be between 15 and 90 minutes", http.StatusBadRequest)
		return
	}

	if preferences.DifficultyPreference < 1 || preferences.DifficultyPreference > 10 {
		http.Error(w, "Difficulty preference must be between 1 and 10", http.StatusBadRequest)
		return
	}

	// Update preferences
	err = h.learningService.UpdateUserPreferences(preferences)
	if err != nil {
		http.Error(w, "Failed to update preferences", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"preferences": preferences,
		"message":     "Preferences updated successfully",
	})
}

// LearningSessionPage renders the immersive learning session page
func (h *LearningHandler) LearningSessionPage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	sessionID := chi.URLParam(r, "sessionID")
	
	// Get user model
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Parse session ID
	parsedSessionID, err := uuid.Parse(sessionID)
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// Get learning session
	session, err := h.learningService.GetLearningSessionByID(parsedSessionID)
	if err != nil {
		http.Error(w, "Failed to get learning session", http.StatusInternalServerError)
		return
	}

	// Verify that the session belongs to the user
	if session.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	// Get the learning prompt for this session
	prompt, err := h.learningService.GetPromptForSession(parsedSessionID)
	if err != nil {
		// Log the error but continue - we'll display the session without the prompt
		fmt.Printf("Warning: failed to get prompt for session: %v\n", err)
	}

	component := templates.LearningSessionPage(user, session, prompt)
	component.Render(r.Context(), w)
}