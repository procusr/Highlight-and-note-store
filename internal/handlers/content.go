package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mh-go/internal/models"
	"mh-go/internal/services"
	"mh-go/templates"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ContentHandler struct {
	contentService *services.ContentService
	folderService  *services.FolderService
}

func NewContentHandler(contentService *services.ContentService, folderService *services.FolderService) *ContentHandler {
	return &ContentHandler{
		contentService: contentService,
		folderService:  folderService,
	}
}

// RegisterRoutes registers all content-related routes
func (h *ContentHandler) RegisterRoutes(r chi.Router, authService *services.AuthService) {
	r.Route("/content", func(r chi.Router) {
		r.Post("/", h.CreateContent)
		r.Get("/{contentID}", h.GetContent)
		r.Put("/{contentID}", h.UpdateContent)
		r.Delete("/{contentID}", h.DeleteContent)
	})
}

// CreateContent creates new content in a folder
func (h *ContentHandler) CreateContent(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uuid.UUID)

	var req struct {
		Title    string    `json:"title"`
		Body     string    `json:"body"`
		FolderID uuid.UUID `json:"folder_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify user owns the folder
	folder, err := h.folderService.GetFolderByID(req.FolderID, userID)
	if err != nil {
		http.Error(w, "Folder not found", http.StatusNotFound)
		return
	}

	if folder.UserID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Create content
	content := &models.Content{
		ID:       uuid.New(),
		Title:    req.Title,
		Body:     req.Body,
		FolderID: req.FolderID,
		UserID:   userID,
	}

	if err := h.contentService.CreateContent(content); err != nil {
		http.Error(w, "Failed to create content", http.StatusInternalServerError)
		return
	}

	// If this is an HTMX request, return the new content card HTML and trigger modal close
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		// Trigger custom event to close the modal
		w.Header().Set("HX-Trigger", "contentCreated")
		if err := templates.ContentCard(*content).Render(r.Context(), w); err != nil {
			http.Error(w, "Error rendering content card", http.StatusInternalServerError)
			return
		}
		return
	}

	// For non-HTMX requests, return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}

// GetContent retrieves content by ID
func (h *ContentHandler) GetContent(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uuid.UUID)
	contentID, err := uuid.Parse(chi.URLParam(r, "contentID"))
	if err != nil {
		http.Error(w, "Invalid content ID", http.StatusBadRequest)
		return
	}

	content, err := h.contentService.GetContentByID(contentID, userID)
	if err != nil {
		http.Error(w, "Content not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}

// UpdateContent updates existing content
func (h *ContentHandler) UpdateContent(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uuid.UUID)
	contentID, err := uuid.Parse(chi.URLParam(r, "contentID"))
	if err != nil {
		http.Error(w, "Invalid content ID", http.StatusBadRequest)
		return
	}

	// Check if user owns the content
	existingContent, err := h.contentService.GetContentByID(contentID, userID)
	if err != nil {
		http.Error(w, "Content not found", http.StatusNotFound)
		return
	}

	var req struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update content
	existingContent.Title = req.Title
	existingContent.Body = req.Body

	if err := h.contentService.UpdateContent(existingContent); err != nil {
		http.Error(w, "Failed to update content", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingContent)
}

// DeleteContent deletes content by ID
func (h *ContentHandler) DeleteContent(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uuid.UUID)
	contentID, err := uuid.Parse(chi.URLParam(r, "contentID"))
	if err != nil {
		http.Error(w, "Invalid content ID", http.StatusBadRequest)
		return
	}

	// Check if user owns the content
	content, err := h.contentService.GetContentByID(contentID, userID)
	if err != nil {
		http.Error(w, "Content not found", http.StatusNotFound)
		return
	}

	if content.UserID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	if err := h.contentService.DeleteContent(contentID, userID); err != nil {
		http.Error(w, "Failed to delete content", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetContentByFolder retrieves all content in a folder with pagination
func (h *ContentHandler) GetContentByFolder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uuid.UUID)
	folderID, err := uuid.Parse(chi.URLParam(r, "folderID"))
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	// Verify user owns the folder
	folder, err := h.folderService.GetFolderByID(folderID, userID)
	if err != nil {
		http.Error(w, "Folder not found", http.StatusNotFound)
		return
	}

	if folder.UserID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Parse pagination parameters
	page := 1
	limit := 10
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	content, err := h.contentService.GetContentByFolderID(folderID, userID)
	if err != nil {
		http.Error(w, "Failed to retrieve content", http.StatusInternalServerError)
		return
	}

	response := struct {
		Content []models.Content `json:"content"`
		Total   int              `json:"total"`
		Page    int              `json:"page"`
		Limit   int              `json:"limit"`
	}{
		Content: content,
		Total:   len(content), // For simplicity, using actual count
		Page:    page,
		Limit:   limit,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
