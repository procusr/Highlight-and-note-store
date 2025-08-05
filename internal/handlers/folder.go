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

		// Content routes
		r.Get("/{folderID}/content", h.GetContent)
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

	// Get user information (simplified for now)
	user := &models.User{ID: userID, Email: "user@example.com"}

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

	// Get user information (simplified for now)
	user := &models.User{ID: userID, Email: "user@example.com"}

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

	// Get user information (simplified for now)
	user := &models.User{ID: userID, Email: "user@example.com"}

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

	// Try form data first
	name := r.FormValue("name")
	description := r.FormValue("description")
	parentIDStr := r.FormValue("parent_id")

	if name != "" {
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
	} else if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
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

	// If this is an HTMX request, redirect to dashboard
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/dashboard")
		w.WriteHeader(http.StatusOK)
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

	// Try form data first
	name := r.FormValue("name")
	description := r.FormValue("description")

	if name != "" {
		req.Name = name
		req.Description = description
	} else if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
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
