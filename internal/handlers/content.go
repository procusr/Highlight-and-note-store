package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mh-go/internal/middleware"
	"mh-go/internal/models"
	"mh-go/internal/services"
	"mh-go/internal/utils"
	"mh-go/templates"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ContentHandler struct {
	contentService *services.ContentService
	folderService  *services.FolderService
	userService    *services.UserService
	openaiService  *services.OpenAIService
}

func NewContentHandler(contentService *services.ContentService, folderService *services.FolderService, userService *services.UserService, openaiService *services.OpenAIService) *ContentHandler {
	return &ContentHandler{
		contentService: contentService,
		folderService:  folderService,
		userService:    userService,
		openaiService:  openaiService,
	}
}

// RegisterRoutes registers all content-related routes
func (h *ContentHandler) RegisterRoutes(r chi.Router, authService *services.AuthService) {
	r.Route("/content", func(r chi.Router) {
		r.Post("/", h.CreateContent)
		r.Get("/{contentID}", h.GetContent)
		r.Get("/{contentID}/view", h.GetViewContentModal)
		r.Get("/{contentID}/edit", h.GetEditContentModal)
		r.Put("/{contentID}", h.UpdateContent)
		r.Delete("/{contentID}", h.DeleteContent)
	})
}

// CreateContent creates new content in a folder
func (h *ContentHandler) CreateContent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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

	// Generate question using OpenAI if the service is available
	if h.openaiService != nil {
		question, err := h.openaiService.GenerateQuestionFromContent(req.Body, req.Title)
		if err == nil {
			content.Question = question
		}
		// If question generation fails, we continue without it - it's not critical
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
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
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

// GetViewContentModal returns the view content modal for a specific content
func (h *ContentHandler) GetViewContentModal(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
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

	w.Header().Set("Content-Type", "text/html")
	if err := templates.ViewContentModal(*content).Render(r.Context(), w); err != nil {
		http.Error(w, "Error rendering view modal", http.StatusInternalServerError)
		return
	}
}

// GetEditContentModal returns the edit content modal HTML
func (h *ContentHandler) GetEditContentModal(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
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

	w.Header().Set("Content-Type", "text/html")
	if err := templates.EditContentModal(*content).Render(r.Context(), w); err != nil {
		http.Error(w, "Error rendering edit modal", http.StatusInternalServerError)
		return
	}
}

// UpdateContent updates existing content
func (h *ContentHandler) UpdateContent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
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

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Update content fields from form
	title := r.FormValue("title")
	description := r.FormValue("description")
	body := r.FormValue("body")

	if title != "" {
		existingContent.Title = title
	}
	if description != "" {
		existingContent.Description = description
	}
	if body != "" {
		existingContent.Body = body
		// Auto-detect content type when body is updated
		existingContent.ContentType = utils.DetectContentType(body)

		// Regenerate question using OpenAI if the service is available
		if h.openaiService != nil {
			question, err := h.openaiService.GenerateQuestionFromContent(body, existingContent.Title)
			if err == nil {
				existingContent.Question = question
			}
			// If question generation fails, we continue without updating it - it's not critical
		}
	}

	if err := h.contentService.UpdateContent(existingContent); err != nil {
		http.Error(w, "Failed to update content", http.StatusInternalServerError)
		return
	}

	// Return the current folder page to show updated content
	folder, err := h.folderService.GetFolderByID(existingContent.FolderID, userID)
	if err != nil {
		http.Error(w, "Failed to load folder", http.StatusInternalServerError)
		return
	}

	// Get user information
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Failed to load user", http.StatusInternalServerError)
		return
	}

	// Get folder contents
	contents, err := h.contentService.GetContentByFolderID(folder.ID, userID)
	if err != nil {
		http.Error(w, "Failed to load contents", http.StatusInternalServerError)
		return
	}

	// Get subfolders
	subfolders, err := h.folderService.GetFoldersByParent(folder.ID, userID)
	if err != nil {
		http.Error(w, "Failed to load subfolders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := templates.FolderPage(user, *folder, contents, subfolders).Render(r.Context(), w); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// DeleteContent deletes content by ID
func (h *ContentHandler) DeleteContent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
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
