package handlers

import (
	"encoding/json"
	"fmt"
	"mh-go/internal/middleware"
	"mh-go/internal/models"
	"mh-go/internal/services"
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

	// Dashboard page route
	r.Get("/dashboard", h.DashboardPage)
}

// DashboardPage renders the main dashboard with folder interface
func (h *FolderHandler) DashboardPage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == uuid.Nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Dashboard - MH-Go</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/htmx.org@1.9.6"></script>
</head>
<body class="bg-gray-50 min-h-screen">
    <div class="flex h-screen">
        <!-- Sidebar -->
        <div class="w-16 bg-gray-800 flex flex-col items-center py-4 space-y-4">
            <div class="w-8 h-8 bg-blue-500 rounded flex items-center justify-center">
                <span class="text-white text-sm">🏠</span>
            </div>
            <div class="w-8 h-8 bg-gray-600 rounded flex items-center justify-center">
                <span class="text-white text-sm">📁</span>
            </div>
            <div class="w-8 h-8 bg-gray-600 rounded flex items-center justify-center">
                <span class="text-white text-sm">📊</span>
            </div>
            <div class="w-8 h-8 bg-gray-600 rounded flex items-center justify-center">
                <span class="text-white text-sm">⚙️</span>
            </div>
            <div class="mt-auto w-8 h-8 bg-gray-600 rounded-full flex items-center justify-center">
                <span class="text-white text-sm">P</span>
            </div>
        </div>

        <!-- Main content -->
        <div class="flex-1 flex flex-col">
            <!-- Header -->
            <div class="bg-white border-b px-6 py-4">
                <div class="flex items-center space-x-4">
                    <!-- Global search -->
                    <div class="flex-1 max-w-md">
                        <input type="text"
                               placeholder="global search"
                               class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                               hx-get="/api/search/content"
                               hx-trigger="keyup changed delay:500ms"
                               hx-target="#search-results">
                    </div>
                    <button class="px-4 py-2 text-gray-600 hover:text-gray-800">Logout</button>
                </div>
            </div>

            <!-- Content area -->
            <div class="flex-1 flex overflow-hidden">
                <!-- Folders panel -->
                <div class="w-80 bg-white border-r flex flex-col">
                    <!-- Folders header -->
                    <div class="border-b px-4 py-3 flex items-center justify-between">
                        <h2 class="font-semibold text-gray-800">Folders</h2>
                        <button class="text-blue-500 hover:text-blue-700"
                                onclick="openCreateFolderModal()">
                            + New
                        </button>
                    </div>

                    <!-- Breadcrumb -->
                    <div id="breadcrumb" class="px-4 py-2 border-b bg-gray-50">
                        <nav class="flex" aria-label="Breadcrumb">
                            <ol class="inline-flex items-center space-x-1 md:space-x-3">
                                <li class="inline-flex items-center">
                                    <a href="#" class="text-gray-700 hover:text-blue-600"
                                       hx-get="/api/folders"
                                       hx-target="#folder-list"
                                       onclick="clearBreadcrumb()">
                                        Root
                                    </a>
                                </li>
                            </ol>
                        </nav>
                    </div>

                    <!-- Folder list -->
                    <div id="folder-list" class="flex-1 overflow-y-auto"
                         hx-get="/api/folders"
                         hx-trigger="load">
                        Loading folders...
                    </div>
                </div>

                <!-- Content panel -->
                <div class="flex-1 flex flex-col">
                    <!-- Content header -->
                    <div class="border-b px-6 py-3 flex items-center justify-between">
                        <h2 id="content-title" class="font-semibold text-gray-800">Select a folder to view contents</h2>
                        <button id="new-content-btn" class="text-blue-500 hover:text-blue-700 hidden"
                                onclick="openCreateContentModal()">
                            + New Content
                        </button>
                    </div>

                    <!-- Content list -->
                    <div id="content-list" class="flex-1 overflow-y-auto p-6">
                        <div class="text-gray-500 text-center">
                            Select a folder from the left panel to view its contents
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- Search results dropdown -->
    <div id="search-results" class="absolute top-16 left-6 right-6 bg-white border rounded-lg shadow-lg z-10 hidden"></div>

    <!-- Modals will go here -->
    <div id="modals"></div>

    <script>
        let currentFolderId = null;

        function selectFolder(folderId, folderName) {
            currentFolderId = folderId;
            document.getElementById('content-title').textContent = folderName + ' Contents';
            document.getElementById('new-content-btn').classList.remove('hidden');

            // Load folder contents
            htmx.ajax('GET', '/api/folders/' + folderId + '/content', '#content-list');

            // Load subfolders
            htmx.ajax('GET', '/api/folders?parent=' + folderId, '#folder-list');

            // Update breadcrumb
            htmx.ajax('GET', '/api/folders/' + folderId + '/path', '#breadcrumb');
        }

        function clearBreadcrumb() {
            currentFolderId = null;
            document.getElementById('content-title').textContent = 'Select a folder to view contents';
            document.getElementById('new-content-btn').classList.add('hidden');
            document.getElementById('content-list').innerHTML = '<div class="text-gray-500 text-center">Select a folder from the left panel to view its contents</div>';
        }

        function openCreateFolderModal() {
            htmx.ajax('GET', '/modals/create-folder?parent=' + (currentFolderId || ''), '#modals');
        }

        function openCreateContentModal() {
            if (currentFolderId) {
                htmx.ajax('GET', '/modals/create-content?folder=' + currentFolderId, '#modals');
            }
        }
    </script>
</body>
</html>`

	w.Write([]byte(html))
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render folder list as HTML
	html := ""
	for _, folder := range folders {
		html += fmt.Sprintf(`
			<div class="p-3 border-b hover:bg-gray-50 cursor-pointer"
			     onclick="selectFolder('%s', '%s')">
				<div class="flex items-center space-x-3">
					<span class="text-yellow-500">📁</span>
					<span class="text-gray-800">%s</span>
				</div>
				<div class="text-sm text-gray-500 ml-6">%s</div>
			</div>`,
			folder.ID, folder.Name, folder.Name, folder.Description)
	}

	if html == "" {
		html = `<div class="p-4 text-gray-500 text-center">No folders found</div>`
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// GetFolder returns a specific folder
func (h *FolderHandler) GetFolder(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	folderID, err := uuid.Parse(chi.URLParam(r, "folderID"))
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	folder, err := h.folderService.GetFolderByID(folderID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folder)
}

// CreateFolder creates a new folder
func (h *FolderHandler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var folder models.Folder
	if err := json.NewDecoder(r.Body).Decode(&folder); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	folder.UserID = userID

	if err := h.folderService.CreateFolder(&folder); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(folder)
}

// UpdateFolder updates an existing folder
func (h *FolderHandler) UpdateFolder(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	folderID, err := uuid.Parse(chi.URLParam(r, "folderID"))
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	var folder models.Folder
	if err := json.NewDecoder(r.Body).Decode(&folder); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	folder.ID = folderID
	folder.UserID = userID

	if err := h.folderService.UpdateFolder(&folder); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folder)
}

// DeleteFolder deletes a folder
func (h *FolderHandler) DeleteFolder(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	folderID, err := uuid.Parse(chi.URLParam(r, "folderID"))
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	if err := h.folderService.DeleteFolder(folderID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetFolderPath returns the breadcrumb path for a folder
func (h *FolderHandler) GetFolderPath(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	folderID, err := uuid.Parse(chi.URLParam(r, "folderID"))
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	path, err := h.folderService.GetFolderPath(folderID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render breadcrumb HTML
	html := `<nav class="flex" aria-label="Breadcrumb">
		<ol class="inline-flex items-center space-x-1 md:space-x-3">
			<li class="inline-flex items-center">
				<a href="#" class="text-gray-700 hover:text-blue-600"
				   hx-get="/api/folders"
				   hx-target="#folder-list"
				   onclick="clearBreadcrumb()">
					Root
				</a>
			</li>`

	for _, item := range path {
		html += fmt.Sprintf(`
			<li>
				<div class="flex items-center">
					<span class="text-gray-400 mx-2">/</span>
					<a href="#" class="text-gray-700 hover:text-blue-600"
					   onclick="selectFolder('%s', '%s')">
						%s
					</a>
				</div>
			</li>`, item.ID, item.Name, item.Name)
	}

	html += `</ol></nav>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// GetContent returns content for a folder
func (h *FolderHandler) GetContent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	folderID, err := uuid.Parse(chi.URLParam(r, "folderID"))
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	contents, err := h.contentService.GetContentByFolderID(folderID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render content list as HTML
	html := ""
	for _, content := range contents {
		html += fmt.Sprintf(`
			<div class="bg-white p-4 rounded-lg border mb-4">
				<div class="flex items-start justify-between">
					<div class="flex-1">
						<h3 class="font-medium text-gray-900">%s</h3>
						<p class="text-sm text-gray-600 mt-1">%s</p>
						<div class="text-xs text-gray-400 mt-2">%s</div>
					</div>
					<div class="flex space-x-2">
						<button class="text-blue-500 hover:text-blue-700 text-sm">Edit</button>
						<button class="text-red-500 hover:text-red-700 text-sm">Delete</button>
					</div>
				</div>
			</div>`,
			content.Title, content.Description, content.CreatedAt.Format("Jan 2, 2006"))
	}

	if html == "" {
		html = `<div class="text-gray-500 text-center">No content found in this folder</div>`
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// CreateContent creates new content in a folder
func (h *FolderHandler) CreateContent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	folderID, err := uuid.Parse(chi.URLParam(r, "folderID"))
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	var content models.Content
	if err := json.NewDecoder(r.Body).Decode(&content); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	content.FolderID = folderID
	content.UserID = userID

	if err := h.contentService.CreateContent(&content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(content)
}

// UpdateContent updates existing content
func (h *FolderHandler) UpdateContent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	folderID, err := uuid.Parse(chi.URLParam(r, "folderID"))
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	contentID, err := uuid.Parse(chi.URLParam(r, "contentID"))
	if err != nil {
		http.Error(w, "Invalid content ID", http.StatusBadRequest)
		return
	}

	var content models.Content
	if err := json.NewDecoder(r.Body).Decode(&content); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	content.ID = contentID
	content.FolderID = folderID
	content.UserID = userID

	if err := h.contentService.UpdateContent(&content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}

// DeleteContent deletes content
func (h *FolderHandler) DeleteContent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	contentID, err := uuid.Parse(chi.URLParam(r, "contentID"))
	if err != nil {
		http.Error(w, "Invalid content ID", http.StatusBadRequest)
		return
	}

	if err := h.contentService.DeleteContent(contentID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SearchContent searches for content
func (h *FolderHandler) SearchContent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	query := r.URL.Query().Get("q")

	if query == "" {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(""))
		return
	}

	contents, err := h.contentService.SearchContent(query, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render search results as HTML
	html := `<div class="max-h-64 overflow-y-auto">`
	for _, content := range contents {
		html += fmt.Sprintf(`
			<div class="p-3 hover:bg-gray-50 border-b">
				<div class="font-medium text-gray-900">%s</div>
				<div class="text-sm text-gray-600">%s</div>
			</div>`,
			content.Title, content.Description)
	}
	html += `</div>`

	if len(contents) == 0 {
		html = `<div class="p-4 text-gray-500 text-center">No results found</div>`
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
