package handlers

import (
	"encoding/json"
	"fmt"
	"mh-go/internal/middleware"
	"mh-go/internal/models"
	"mh-go/internal/services"
	"mh-go/templates"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type FolderHandler struct {
	folderService  *services.FolderService
	contentService *services.ContentService
}

func NewFolderHandler(folderService *services.FolderService, contentService *services.ContentService) *FolderHandler {
	return &FolderHandler{
		folderService:  folderService,
		contentService: contentService,
	}
}

// RegisterRoutes registers all folder-related routes
func (h *FolderHandler) RegisterRoutes(r chi.Router, authService *services.AuthService) {
	r.Route("/api/folders", func(r chi.Router) {
		r.Use(middleware.RequireAuth(authService))

		// Folder routes
		r.Get("/", h.GetFolders)
		r.Post("/", h.CreateFolder)
		r.Get("/{folderID}", h.GetFolder)
		r.Put("/{folderID}", h.UpdateFolder)
		r.Delete("/{folderID}", h.DeleteFolder)
		r.Get("/{folderID}/path", h.GetFolderPath)

		// Content routes
		r.Get("/{folderID}/content", h.GetContent)
		r.Post("/{folderID}/content", h.CreateContent)
		r.Put("/{folderID}/content/{contentID}", h.UpdateContent)
		r.Delete("/{folderID}/content/{contentID}", h.DeleteContent)
	})

	// Search route
	r.Route("/api/search", func(r chi.Router) {
		r.Use(middleware.RequireAuth(authService))
		r.Get("/content", h.SearchContent)
	})

	// Page routes
	r.Get("/dashboard", h.DashboardPage)
	r.Get("/folders", h.FoldersPage)
	r.Get("/folders/{folderID}", h.FolderPage)
}

// DashboardPage renders the main dashboard with folder interface
func (h *FolderHandler) DashboardPage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get user information
	user := &models.User{ID: userID} // We should fetch this from user service, but for now we'll use the ID

	// Get root folders for the user
	folders, err := h.folderService.GetRootFolders(userID)
	if err != nil {
		http.Error(w, "Failed to load folders", http.StatusInternalServerError)
		return
	}

	// Render dashboard template
	w.Header().Set("Content-Type", "text/html")
	if err := templates.DashboardPage(user, folders).Render(r.Context(), w); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// FoldersPage renders the folders list page
func (h *FolderHandler) FoldersPage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get user information
	user := &models.User{ID: userID} // We should fetch this from user service, but for now we'll use the ID

	// Get all folders for the user
	folders, err := h.folderService.GetRootFolders(userID)
	if err != nil {
		http.Error(w, "Failed to load folders", http.StatusInternalServerError)
		return
	}

	// Render folders template
	w.Header().Set("Content-Type", "text/html")
	if err := templates.FoldersPage(user, folders).Render(r.Context(), w); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// FolderPage renders a specific folder view
func (h *FolderHandler) FolderPage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	folderIDStr := chi.URLParam(r, "folderID")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	// Get folder details
	folder, err := h.folderService.GetFolderByID(folderID, userID)
	if err != nil {
		http.Error(w, "Folder not found", http.StatusNotFound)
		return
	}

	// Get folder contents
	contents, err := h.contentService.GetContentByFolderID(folderID, userID)
	if err != nil {
		http.Error(w, "Failed to load folder contents", http.StatusInternalServerError)
		return
	}

	// Get subfolders
	subfolders, err := h.folderService.GetFoldersByParent(folderID, userID)
	if err != nil {
		http.Error(w, "Failed to load subfolders", http.StatusInternalServerError)
		return
	}

	// Get user information
	user := &models.User{ID: userID} // We should fetch this from user service, but for now we'll use the ID

	// Render folder page template
	w.Header().Set("Content-Type", "text/html")
	if err := templates.FolderPage(user, *folder, contents, subfolders).Render(r.Context(), w); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// GetFolders returns folders (root folders or children of a parent)
func (h *FolderHandler) GetFolders(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	parentID := r.URL.Query().Get("parent")

	var folders []models.Folder
	var err error

	if parentID == "" {
		folders, err = h.folderService.GetRootFolders(userID)
	} else {
		parentUUID, parseErr := uuid.Parse(parentID)
		if parseErr != nil {
			http.Error(w, "Invalid parent ID", http.StatusBadRequest)
			return
		}
		folders, err = h.folderService.GetFoldersByParent(parentUUID, userID)
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get folders: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folders)
}

// CreateFolder creates a new folder
func (h *FolderHandler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req struct {
		Name        string     `json:"name"`
		Description string     `json:"description"`
		ParentID    *uuid.UUID `json:"parent_id"`
	}

	// Check content type to determine how to parse the request
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}
	} else {
		// Parse as form data (default for HTMX forms)
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		description := r.FormValue("description")
		parentIDStr := r.FormValue("parent_id")

		if name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}

		req.Name = name
		req.Description = description

		if parentIDStr != "" {
			parentID, err := uuid.Parse(parentIDStr)
			if err != nil {
				http.Error(w, "Invalid parent ID", http.StatusBadRequest)
				return
			}
			req.ParentID = &parentID
		}
	}

	folder := &models.Folder{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		UserID:      userID,
	}

	if err := h.folderService.CreateFolder(folder); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create folder: %v", err), http.StatusInternalServerError)
		return
	}

	// If this is an HTMX request, return the folder card
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		if err := templates.FolderCard(*folder).Render(r.Context(), w); err != nil {
			http.Error(w, "Error rendering folder card", http.StatusInternalServerError)
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(folder)
}

// GetFolder returns a specific folder
func (h *FolderHandler) GetFolder(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	folderIDStr := chi.URLParam(r, "folderID")

	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	folder, err := h.folderService.GetFolderByID(folderID, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get folder: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folder)
}

// UpdateFolder updates an existing folder
func (h *FolderHandler) UpdateFolder(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	folderIDStr := chi.URLParam(r, "folderID")

	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	// Check if folder exists and user owns it
	existingFolder, err := h.folderService.GetFolderByID(folderID, userID)
	if err != nil {
		http.Error(w, "Folder not found", http.StatusNotFound)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Try form data
		name := r.FormValue("name")
		description := r.FormValue("description")

		if name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}

		req.Name = name
		req.Description = description
	}

	existingFolder.Name = req.Name
	existingFolder.Description = req.Description

	if err := h.folderService.UpdateFolder(existingFolder); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update folder: %v", err), http.StatusInternalServerError)
		return
	}

	// If this is an HTMX request, redirect to dashboard
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/dashboard")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingFolder)
}

// DeleteFolder deletes a folder
func (h *FolderHandler) DeleteFolder(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	folderIDStr := chi.URLParam(r, "folderID")

	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	// Check if folder exists and user owns it
	folder, err := h.folderService.GetFolderByID(folderID, userID)
	if err != nil {
		http.Error(w, "Folder not found", http.StatusNotFound)
		return
	}

	if folder.UserID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	if err := h.folderService.DeleteFolder(folderID, userID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete folder: %v", err), http.StatusInternalServerError)
		return
	}

	// If this is an HTMX request, redirect to dashboard
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/dashboard")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetContent returns content for a folder
func (h *FolderHandler) GetContent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	folderIDStr := chi.URLParam(r, "folderID")

	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	content, err := h.contentService.GetContentByFolderID(folderID, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get content: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}

// CreateContent creates new content in a folder
func (h *FolderHandler) CreateContent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	folderIDStr := chi.URLParam(r, "folderID")

	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Body        string `json:"body"`
		ContentType string `json:"content_type"`
	}

	// Check content type to determine how to parse the request
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}
	} else {
		// Parse as form data (default for HTMX forms)
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		description := r.FormValue("description")
		body := r.FormValue("body")
		contentTypeForm := r.FormValue("content_type")

		if title == "" || body == "" {
			http.Error(w, "Title and body are required", http.StatusBadRequest)
			return
		}

		req.Title = title
		req.Description = description
		req.Body = body
		req.ContentType = contentTypeForm
	}

	if req.ContentType == "" {
		req.ContentType = "text"
	}

	content := &models.Content{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		Body:        req.Body,
		ContentType: req.ContentType,
		FolderID:    folderID,
		UserID:      userID,
	}

	if err := h.contentService.CreateContent(content); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create content: %v", err), http.StatusInternalServerError)
		return
	}

	// If this is an HTMX request, redirect to folder page
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", fmt.Sprintf("/folders/%s", folderID.String()))
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(content)
}

// UpdateContent updates existing content
func (h *FolderHandler) UpdateContent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	folderIDStr := chi.URLParam(r, "folderID")
	contentIDStr := chi.URLParam(r, "contentID")

	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	contentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		http.Error(w, "Invalid content ID", http.StatusBadRequest)
		return
	}

	// Check if content exists and user owns it
	existingContent, err := h.contentService.GetContentByID(contentID, userID)
	if err != nil {
		http.Error(w, "Content not found", http.StatusNotFound)
		return
	}

	if existingContent.FolderID != folderID {
		http.Error(w, "Content not in specified folder", http.StatusBadRequest)
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Body        string `json:"body"`
		ContentType string `json:"content_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Try form data
		title := r.FormValue("title")
		description := r.FormValue("description")
		body := r.FormValue("body")
		contentType := r.FormValue("content_type")

		if title == "" || body == "" {
			http.Error(w, "Title and body are required", http.StatusBadRequest)
			return
		}

		req.Title = title
		req.Description = description
		req.Body = body
		req.ContentType = contentType
	}

	existingContent.Title = req.Title
	existingContent.Description = req.Description
	existingContent.Body = req.Body
	if req.ContentType != "" {
		existingContent.ContentType = req.ContentType
	}

	if err := h.contentService.UpdateContent(existingContent); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update content: %v", err), http.StatusInternalServerError)
		return
	}

	// If this is an HTMX request, redirect to folder page
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", fmt.Sprintf("/folders/%s", folderID.String()))
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingContent)
}

// DeleteContent deletes content from a folder
func (h *FolderHandler) DeleteContent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	folderIDStr := chi.URLParam(r, "folderID")
	contentIDStr := chi.URLParam(r, "contentID")

	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	contentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		http.Error(w, "Invalid content ID", http.StatusBadRequest)
		return
	}

	// Check if content exists and user owns it
	content, err := h.contentService.GetContentByID(contentID, userID)
	if err != nil {
		http.Error(w, "Content not found", http.StatusNotFound)
		return
	}

	if content.FolderID != folderID {
		http.Error(w, "Content not in specified folder", http.StatusBadRequest)
		return
	}

	if content.UserID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	if err := h.contentService.DeleteContent(contentID, userID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete content: %v", err), http.StatusInternalServerError)
		return
	}

	// If this is an HTMX request, redirect to folder page
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", fmt.Sprintf("/folders/%s", folderID.String()))
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SearchContent searches for content across all folders
func (h *FolderHandler) SearchContent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	query := r.URL.Query().Get("q")

	if query == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]models.Content{})
		return
	}

	content, err := h.contentService.SearchContent(query, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to search content: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}

// GetFolderPath returns the full path of a folder
func (h *FolderHandler) GetFolderPath(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	folderIDStr := chi.URLParam(r, "folderID")

	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	path, err := h.folderService.GetFolderPath(folderID, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get folder path: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"path": path,
	})
}
