package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"mh-go/internal/middleware"
	"mh-go/internal/models"
	"mh-go/internal/services"
	"mh-go/internal/utils"
	"mh-go/templates"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type FolderHandler struct {
	folderService  *services.FolderService
	contentService *services.ContentService
	userService    *services.UserService
	openaiService  *services.OpenAIService
}

func NewFolderHandler(folderService *services.FolderService, contentService *services.ContentService, userService *services.UserService, openaiService *services.OpenAIService) *FolderHandler {
	return &FolderHandler{
		folderService:  folderService,
		contentService: contentService,
		userService:    userService,
		openaiService:  openaiService,
	}
}

// RegisterRoutes registers all folder-related routes
func (h *FolderHandler) RegisterRoutes(r chi.Router, authService *services.AuthService) {
	r.Route("/folders", func(r chi.Router) {
		r.Use(middleware.RequireAuth(authService))

		// Folder routes
		r.Get("/", h.GetFolders)
		r.Post("/", h.CreateFolder)
		r.Get("/{folderID}", h.GetFolder)
		r.Get("/{folderID}/edit", h.GetEditFolderModal)
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
	r.Route("/search", func(r chi.Router) {
		r.Use(middleware.RequireAuth(authService))
		r.Get("/content", h.SearchContent)
	})
}

// OverviewPage renders the main overview with learning progress interface
func (h *FolderHandler) OverviewPage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get user information
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Failed to load user", http.StatusInternalServerError)
		return
	}

	// Render overview template
	w.Header().Set("Content-Type", "text/html")
	if err := templates.OverviewPage(user).Render(r.Context(), w); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// FoldersPage renders the folders listing page
func (h *FolderHandler) FoldersPage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get user information
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Failed to load user", http.StatusInternalServerError)
		return
	}

	// Get all folders for the user (including subfolders)
	folders, err := h.folderService.GetRootFolders(userID)
	if err != nil {
		http.Error(w, "Failed to load folders", http.StatusInternalServerError)
		return
	}

	// Render folders page template
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

	// Get user details
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Failed to load user", http.StatusInternalServerError)
		return
	}

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
	if userID == uuid.Nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name        string     `json:"name"`
		Description string     `json:"description"`
		ParentID    *uuid.UUID `json:"parent_id"`
	}

	contentType := r.Header.Get("Content-Type")
	fmt.Printf("Content-Type: %s\n", contentType)

	if contentType == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			fmt.Printf("JSON decode error: %v\n", err)
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}
	} else {
		// Parse form data first
		if parseErr := r.ParseForm(); parseErr != nil {
			fmt.Printf("Error parsing form: %v\n", parseErr)
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		// Try form data
		name := r.FormValue("name")
		description := r.FormValue("description")
		parentIDStr := r.FormValue("parent_id")

		fmt.Printf("Form data: name='%s', description='%s', parent_id='%s'\n", name, description, parentIDStr)
		fmt.Printf("All form values: %+v\n", r.Form)

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

	// If this is an HTMX request, return the new folder card HTML
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		// Trigger custom event to close the modal
		w.Header().Set("HX-Trigger", "folderCreated")
		if err := templates.FolderCard(*folder).Render(r.Context(), w); err != nil {
			http.Error(w, "Error rendering folder card", http.StatusInternalServerError)
			return
		}
		return
	}

	// For non-HTMX requests, return JSON
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

// GetEditFolderModal returns the edit folder modal HTML
func (h *FolderHandler) GetEditFolderModal(w http.ResponseWriter, r *http.Request) {
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

	// Check if folder can be edited
	if !folder.IsEditable() {
		http.Error(w, "This folder cannot be edited", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := templates.EditFolderModal(*folder).Render(r.Context(), w); err != nil {
		http.Error(w, "Error rendering edit modal", http.StatusInternalServerError)
		return
	}
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

	// If this is an HTMX request, redirect to folders page
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/folders")
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

	// If this is an HTMX request, redirect to folders page
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/folders")
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
	if userID == uuid.Nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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
	}

	contentType := r.Header.Get("Content-Type")
	fmt.Printf("Content creation - Content-Type: %s\n", contentType)

	if contentType == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			fmt.Printf("Content creation - JSON decode error: %v\n", err)
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}
	} else {
		// Parse form data first
		if parseErr := r.ParseForm(); parseErr != nil {
			fmt.Printf("Content creation - Error parsing form: %v\n", parseErr)
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		// Try form data
		title := r.FormValue("title")
		description := r.FormValue("description")
		body := r.FormValue("body")

		fmt.Printf("Content form data: title='%s', description='%s', body='%s'\n", title, description, body)

		if body == "" {
			http.Error(w, "Content body is required", http.StatusBadRequest)
			return
		}

		req.Title = title
		req.Description = description
		req.Body = body
	}

	// Always attempt to generate a question (primary heading) if body present
	var generatedQuestion string
	if strings.TrimSpace(req.Body) == "" {
		log.Println("[CreateContent] Empty body; skipping question generation")
	} else {
		if h.openaiService != nil {
			log.Println("[CreateContent] Attempting AI question generation...")
			if question, err := h.openaiService.GenerateQuestionFromContent(req.Body, req.Title); err != nil {
				log.Printf("[CreateContent] AI generation failed: %v -- using heuristic fallback", err)
				generatedQuestion = heuristicQuestion(req.Body, req.Title)
			} else {
				generatedQuestion = strings.TrimSpace(question)
				if generatedQuestion == "" {
					log.Println("[CreateContent] AI returned empty content; using heuristic fallback")
					generatedQuestion = heuristicQuestion(req.Body, req.Title)
				} else {
					log.Printf("[CreateContent] AI question generated (%d chars)", len(generatedQuestion))
				}
			}
		} else {
			log.Println("[CreateContent] OpenAI service not configured; using heuristic question")
			generatedQuestion = heuristicQuestion(req.Body, req.Title)
		}
	}

	// If user supplied a title but no question produced (unlikely), still fallback heuristic
	if generatedQuestion == "" {
		generatedQuestion = heuristicQuestion(req.Body, req.Title)
	}

	// If user provided no title, keep title empty (we show question). If both absent we fallback for title just for storage/search.
	if strings.TrimSpace(req.Title) == "" {
		fallbackTitle := deriveFallbackTitle(req.Body)
		if fallbackTitle != "" {
			req.Title = fallbackTitle
		}
	}

	// Automatically detect content type based on content
	detectedContentType := utils.DetectContentType(req.Body)

	content := &models.Content{
		ID:          uuid.New(),
		Title:       req.Title,
		Question:    generatedQuestion,
		Description: req.Description,
		Body:        req.Body,
		ContentType: detectedContentType,
		FolderID:    folderID,
		UserID:      userID,
	}

	if err := h.contentService.CreateContent(content); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create content: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("[CreateContent] Saved content ID=%s Title=%q QuestionPresent=%t", content.ID, content.Title, generatedQuestion != "")

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

// deriveFallbackTitle builds a short heading from the body when neither title nor AI question is available.
// Strategy:
// 1. Split into lines, pick first non-empty trimmed line.
// 2. If still empty, take first 80 runes of full body.
// 3. Collapse spaces, trim punctuation, ensure <= 80 chars, append ellipsis if truncated.
func deriveFallbackTitle(body string) string {
	trimmed := strings.TrimSpace(body)
	if trimmed == "" {
		return ""
	}

	// Prefer first meaningful line
	lines := strings.Split(trimmed, "\n")
	var candidate string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			candidate = l
			break
		}
	}
	if candidate == "" {
		candidate = trimmed
	}

	// Normalize whitespace
	candidate = strings.Join(strings.Fields(candidate), " ")

	// Truncate to 80 runes
	const limit = 80
	if utf8.RuneCountInString(candidate) > limit {
		// build rune slice safely
		runes := []rune(candidate)
		candidate = string(runes[:limit]) + "…"
	}

	// Remove trailing punctuation clusters
	candidate = strings.TrimRight(candidate, ":;,. -")

	return candidate
}

// heuristicQuestion creates a simple question when AI generation is unavailable.
// Strategy: take first meaningful sentence or line; transform to a how/what/why style question.
func heuristicQuestion(body, title string) string {
	// Prefer explicit title if present to craft question
	if t := strings.TrimSpace(title); t != "" {
		return ensureQuestion(t)
	}
	body = strings.TrimSpace(body)
	if body == "" {
		return "What is the main idea?"
	}

	// Extract first sentence-like fragment
	first := body
	// Split by sentence end punctuation
	if idx := regexp.MustCompile(`[.!?]`).FindStringIndex(body); idx != nil && idx[0] > 15 { // require a little length
		first = body[:idx[0]]
	} else {
		// fallback first line
		if nl := strings.Index(body, "\n"); nl > 0 {
			first = body[:nl]
		}
	}
	first = strings.TrimSpace(first)
	if first == "" {
		return "What is the main idea?"
	}
	// Remove leading markdown symbols
	first = strings.TrimLeft(first, "#*- ")
	// Truncate
	runes := []rune(first)
	if len(runes) > 90 {
		first = string(runes[:90])
	}
	return ensureQuestion(first)
}

// ensureQuestion ensures text ends with a question mark and starts with an interrogative if reasonable.
func ensureQuestion(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "What is the main idea?"
	}
	// If already ends with ? assume it's a question
	if strings.HasSuffix(s, "?") {
		return s
	}
	lower := strings.ToLower(s)
	startsQuestion := false
	for _, w := range []string{"what", "why", "how", "when", "where", "who", "which"} {
		if strings.HasPrefix(lower, w+" ") {
			startsQuestion = true
			break
		}
	}
	if !startsQuestion {
		// Prepend "What is" or "How does" heuristically based on length
		if len(strings.Fields(s)) <= 4 {
			s = "What is " + s
		} else {
			s = "How does " + s
		}
	}
	return strings.TrimSpace(s) + "?"
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
