package handlers

import (
	"encoding/json"
	"fmt"
	"mh-go/internal/middleware"
	"mh-go/internal/models"
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

// getStaticPathwayCards returns the hardcoded learning pathway cards
func (h *LearningHandler) getStaticPathwayCards() []models.LearningPathwayCard {
	return []models.LearningPathwayCard{
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Name:        "Applied Mastery Laboratory",
			Description: "Hands-on practice and real-world application of concepts through experiments and projects.",
			MethodKey:   "applied_laboratory",
			HasStarted:  false,
			ProgressPercent: 0,
		},
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Name:        "Foundation Through Discovery",
			Description: "Explore and discover fundamental concepts through guided inquiry and investigation.",
			MethodKey:   "foundation_discovery",
			HasStarted:  false,
			ProgressPercent: 0,
		},
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			Name:        "Integration & Transfer System",
			Description: "Connect knowledge across domains and apply learning to new contexts and situations.",
			MethodKey:   "integration_transfer",
			HasStarted:  false,
			ProgressPercent: 0,
		},
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			Name:        "Mistake-Based Learning Lab",
			Description: "Learn from errors and develop resilience through controlled failure analysis and correction.",
			MethodKey:   "mistake_learning_lab",
			HasStarted:  false,
			ProgressPercent: 0,
		},
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000000005"),
			Name:        "Problem-Solving Gymnasium",
			Description: "Build critical thinking skills through structured problem-solving exercises and challenges.",
			MethodKey:   "problem_gymnasium",
			HasStarted:  false,
			ProgressPercent: 0,
		},
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
		// Create default stats if error
		stats = &models.DashboardStats{
			OverallProgress:   0,
			ConceptsMastered: 0,
			DailyStreak:       0,
			TotalStudyTime:    0,
		}
	}

	// Get learning pathway cards (static data)
	pathwayCards := h.getStaticPathwayCards()

	// Get available topics
	topics, err := h.topicService.GetAllTopics()
	if err != nil {
		topics = []models.Topic{}
	}

	// Get active learning sessions
	sessions, err := h.learningService.GetActiveLearningSessions(userID)
	if err != nil {
		sessions = []models.LearningSession{}
	}
	fmt.Printf("[LearningPlatformPage] user=%s active_sessions=%d\n", userID.String(), len(sessions))

	// Get user preferences
	preferences, err := h.learningService.GetOrCreateUserPreferences(userID)
	if err != nil {
		// Create default preferences
		preferences = &models.UserPreferences{
			UserID:             userID,
			LearningIntensity:  5,
			SessionDurationMins: 30,
			DifficultyPreference: 5,
			ThemePreference:    "light",
		}
	}

	component := templates.LearningPlatformPage(user, stats, pathwayCards, topics, sessions, preferences)
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

	topics, _, err := h.learningService.SearchTopics(query, 1, limit)
	if err != nil {
		http.Error(w, "Failed to search topics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topics)
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
		Topic        string `json:"topic"`
		MethodKey    string `json:"method_key"`
		LearningGoal string `json:"learning_goal"`
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "" && (contentType == "application/json" || contentType[:16] == "application/json") {
		// JSON body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	} else {
		// Assume form submission (e.g., from HTMX form post)
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form submission", http.StatusBadRequest)
			return
		}
		req.Topic = r.FormValue("topic")
		req.MethodKey = r.FormValue("method_key")
		req.LearningGoal = r.FormValue("learning_goal")
	}

	// Create a topic for the user's learning session
	topic := &models.Topic{
		Name:            req.Topic,
		Description:     fmt.Sprintf("Learning topic: %s", req.Topic),
		Category:        "user-generated",
		DifficultyLevel: "intermediate",
	}

	err := h.topicService.CreateTopic(topic)
	if err != nil {
		http.Error(w, "Failed to create topic", http.StatusInternalServerError)
		return
	}
	topicID := topic.ID

	// Get learning method by key
	method, err := h.methodService.GetLearningMethodByKey(req.MethodKey)
	if err != nil {
		http.Error(w, "Invalid learning method", http.StatusBadRequest)
		return
	}

	// Create learning prompt
	_, err = h.learningService.CreateLearningPrompt(userID, topicID, method.ID, req.LearningGoal)
	if err != nil {
		http.Error(w, "Failed to create learning prompt", http.StatusInternalServerError)
		return
	}

	// Start learning session
	session, err := h.learningService.StartLearningSession(userID, topicID, method.ID)
	if err != nil {
		http.Error(w, "Failed to start learning session", http.StatusInternalServerError)
		return
	}
	fmt.Printf("[StartLearningSession] user=%s session=%s topic=%s method=%s method_key=%s\n", userID.String(), session.ID.String(), topicID.String(), method.ID.String(), method.MethodKey)

	// If HTMX request, issue redirect header for refresh
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/learning")
		w.WriteHeader(http.StatusNoContent)
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
	sessionID, err := uuid.Parse(chi.URLParam(r, "sessionID"))
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// Get learning session
	session, err := h.learningService.GetLearningSessionByID(sessionID)
	if err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Verify the session belongs to the current user
	if session.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the learning prompt for this session
	prompt, err := h.learningService.GetPromptForSession(sessionID)
	if err != nil {
		// It's okay if there's no prompt yet
		prompt = nil
	}

	// Get user model
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	component := templates.LearningSessionPage(user, session, prompt)
	component.Render(r.Context(), w)
}

// DeleteLearningSession deletes a learning session
func (h *LearningHandler) DeleteLearningSession(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	sessionID, err := uuid.Parse(chi.URLParam(r, "sessionID"))
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	err = h.learningService.DeleteLearningSession(sessionID, userID)
	if err != nil {
		if err.Error() == "session not found" {
			http.Error(w, "Session not found", http.StatusNotFound)
			return
		}
		if err.Error() == "unauthorized: session does not belong to user" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Failed to delete learning session", http.StatusInternalServerError)
		return
	}

	// If HTMX request, return success response
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Learning session deleted successfully",
	})
}