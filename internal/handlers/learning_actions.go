package handlers

import (
	"fmt"
	"mh-go/internal/middleware"
	"mh-go/internal/models"
	"mh-go/templates"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// StartLearning handles the form submission to start a learning session
func (h *LearningHandler) StartLearning(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	methodID := strings.TrimSpace(r.FormValue("method_id"))
	methodKey := strings.TrimSpace(r.FormValue("method_key"))
	topic := strings.TrimSpace(r.FormValue("topic"))
	learningGoal := strings.TrimSpace(r.FormValue("learning_goal"))

	if topic == "" || learningGoal == "" || (methodID == "" && methodKey == "") {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`
			<div class="bg-red-50 border border-red-200 rounded-md p-4">
				<div class="flex">
					<div class="ml-3">
						<h3 class="text-sm font-medium text-red-800">Error</h3>
						<div class="mt-2 text-sm text-red-700">
							<p>All fields are required (topic, goal, method).</p>
						</div>
					</div>
				</div>
			</div>
		`))
		return
	}

	// Resolve learning method
	var resolvedMethodID uuid.UUID
	var resolvedMethodKey string
	if methodID != "" {
		parsed, err := uuid.Parse(methodID)
		if err != nil {
			http.Error(w, "Invalid method ID", http.StatusBadRequest)
			return
		}
		resolvedMethodID = parsed
		// Try load method to get key if not provided
		if methodKey == "" {
			if m, err := h.methodService.GetLearningMethodByID(parsed); err == nil {
				resolvedMethodKey = m.MethodKey
			}
		} else {
			resolvedMethodKey = methodKey
		}
	} else {
		// Use method key lookup
		m, err := h.methodService.GetLearningMethodByKey(methodKey)
		if err != nil {
			http.Error(w, "Unknown learning method", http.StatusBadRequest)
			return
		}
		resolvedMethodID = m.ID
		resolvedMethodKey = m.MethodKey
	}

	// Find or create topic
	var topicModel *models.Topic
	topics, _, err := h.learningService.SearchTopics(topic, 1, 1)
	if err == nil && len(topics) > 0 {
		topicModel = &topics[0]
	} else {
		topicModel = &models.Topic{
			Name:            topic,
			Description:     fmt.Sprintf("Custom topic: %s", topic),
			Category:        "Custom",
			DifficultyLevel: "intermediate",
		}
		if err := h.topicService.CreateTopic(topicModel); err != nil {
			http.Error(w, "Failed to create topic", http.StatusInternalServerError)
			return
		}
	}

	// Create learning prompt
	prompt, err := h.learningService.CreateLearningPrompt(userID, topicModel.ID, resolvedMethodID, learningGoal)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`
			<div class="bg-red-50 border border-red-200 rounded-md p-4">
				<div class="flex">
					<div class="ml-3">
						<h3 class="text-sm font-medium text-red-800">Error</h3>
						<div class="mt-2 text-sm text-red-700">
							<p>Failed to create learning session (prompt error). Please try again.</p>
						</div>
					</div>
				</div>
			</div>
		`))
		fmt.Printf("[StartLearning] ERROR prompt creation user=%s topic=%s method_key=%s err=%v\n", userID, topic, resolvedMethodKey, err)
		return
	}

	session, err := h.learningService.StartLearningSession(userID, topicModel.ID, resolvedMethodID)
	if err != nil {
		fmt.Printf("[StartLearning] WARNING failed to create session user=%s topic=%s method_key=%s err=%v\n", userID, topic, resolvedMethodKey, err)
	}

	if session != nil {
		fmt.Printf("[StartLearning] SUCCESS user=%s topic=%s method_key=%s method_id=%s prompt_id=%s session_id=%s\n",
			userID, topic, resolvedMethodKey, resolvedMethodID, prompt.ID, session.ID)
	} else {
		fmt.Printf("[StartLearning] SUCCESS (no session) user=%s topic=%s method_key=%s method_id=%s prompt_id=%s\n",
			userID, topic, resolvedMethodKey, resolvedMethodID, prompt.ID)
	}

	// HTMX redirect to refresh page so new session appears
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/learning")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Fallback: render AI prompt component (non-HTMX direct post)
	component := templates.AIPromptDisplay(prompt, topic, resolvedMethodKey)
	component.Render(r.Context(), w)
}

