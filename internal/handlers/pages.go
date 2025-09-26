package handlers

import (
	"mh-go/internal/middleware"
	"mh-go/internal/models"
	"mh-go/internal/services"
	"mh-go/templates"
	"net/http"
)

type PageHandler struct {
	userService *services.UserService
}

func NewPageHandler(userService *services.UserService) *PageHandler {
	return &PageHandler{
		userService: userService,
	}
}

// FlashcardsPage renders the flashcards page
func (h *PageHandler) FlashcardsPage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	component := templates.FlashcardsPage(user)
	component.Render(r.Context(), w)
}

// RememberPage renders the remember (spaced repetition) page
func (h *PageHandler) RememberPage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	component := templates.RememberPage(userID)
	component.Render(r.Context(), w)
}

// AnalyzePage renders the analytics page
func (h *PageHandler) AnalyzePage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	component := templates.AnalyzePage(userID)
	component.Render(r.Context(), w)
}

// OptimizePage renders the learning optimization page
func (h *PageHandler) OptimizePage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	component := templates.OptimizePage(userID)
	component.Render(r.Context(), w)
}

// SettingsPage renders the settings page
func (h *PageHandler) SettingsPage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	component := templates.SettingsPage(user)
	component.Render(r.Context(), w)
}
