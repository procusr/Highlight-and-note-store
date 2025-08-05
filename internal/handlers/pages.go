package handlers

import (
	"mh-go/internal/models"
	"mh-go/templates"
	"net/http"
)

type PageHandler struct{}

func NewPageHandler() *PageHandler {
	return &PageHandler{}
}

// FlashcardsPage renders the flashcards page
func (h *PageHandler) FlashcardsPage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	component := templates.FlashcardsPage(user)
	component.Render(r.Context(), w)
}

// SettingsPage renders the settings page
func (h *PageHandler) SettingsPage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	component := templates.SettingsPage(user)
	component.Render(r.Context(), w)
}
