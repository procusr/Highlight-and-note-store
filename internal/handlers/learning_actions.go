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

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	methodID := r.FormValue("method_id")
	methodKey := r.FormValue("method_key")
	topic := strings.TrimSpace(r.FormValue("topic"))
	userPrompt := strings.TrimSpace(r.FormValue("user_prompt"))

	if methodID == "" || topic == "" || userPrompt == "" {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`
			<div class="bg-red-50 border border-red-200 rounded-md p-4">
				<div class="flex">
					<div class="ml-3">
						<h3 class="text-sm font-medium text-red-800">Error</h3>
						<div class="mt-2 text-sm text-red-700">
							<p>All fields are required. Please fill in the topic and learning goal.</p>
						</div>
					</div>
				</div>
			</div>
		`))
		return
	}

	// Parse method ID
	parsedMethodID, err := uuid.Parse(methodID)
	if err != nil {
		http.Error(w, "Invalid method ID", http.StatusBadRequest)
		return
	}

	// Find or create topic
	var topicModel *models.Topic
	
	// First try to find existing topic by name
	topics, err := h.learningService.SearchTopics(topic, 1)
	if err == nil && len(topics) > 0 {
		topicModel = &topics[0]
	} else {
		// Create new topic
		topicModel = &models.Topic{
			Name:            topic,
			Description:     fmt.Sprintf("Custom topic: %s", topic),
			Category:        "Custom",
			DifficultyLevel: "intermediate",
		}
		
		err = h.topicService.CreateTopic(topicModel)
		if err != nil {
			http.Error(w, "Failed to create topic", http.StatusInternalServerError)
			return
		}
	}

	// Create learning prompt
	prompt, err := h.learningService.CreateLearningPrompt(userID, topicModel.ID, parsedMethodID, userPrompt)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`
			<div class="bg-red-50 border border-red-200 rounded-md p-4">
				<div class="flex">
					<div class="ml-3">
						<h3 class="text-sm font-medium text-red-800">Error</h3>
						<div class="mt-2 text-sm text-red-700">
							<p>Failed to create learning session. Please try again.</p>
						</div>
					</div>
				</div>
			</div>
		`))
		return
	}

	// Generate AI prompt display
	component := templates.AIPromptDisplay(prompt, topic, methodKey)
	component.Render(r.Context(), w)
}

