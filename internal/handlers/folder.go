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

	// Modal routes
	r.Route("/modals", func(r chi.Router) {
		r.Use(middleware.RequireAuth(authService))
		r.Get("/create-folder", h.CreateFolderModal)
		r.Get("/edit-folder/{folderID}", h.EditFolderModal)
		r.Get("/create-content", h.CreateContentModal)
		r.Get("/edit-content/{contentID}", h.EditContentModal)
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
    <style>
        .mobile-nav { display: none; }
        @media (max-width: 768px) {
            .desktop-sidebar { display: none; }
            .mobile-nav { display: flex; }
            .folder-panel { width: 100%; }
            .content-panel { width: 100%; }
        }
        .loading {
            opacity: 0.6;
            pointer-events: none;
        }
        .fade-in {
            animation: fadeIn 0.3s ease-in-out;
        }
        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(10px); }
            to { opacity: 1; transform: translateY(0); }
        }
    </style>
</head>
<body class="bg-gray-50 min-h-screen">
    <!-- Mobile Navigation -->
    <div class="mobile-nav bg-white border-b px-4 py-3 flex items-center justify-between md:hidden">
        <button id="mobile-menu-btn" class="p-2 rounded-md hover:bg-gray-100">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
            </svg>
        </button>
        <h1 class="text-lg font-semibold">MH-Go</h1>
        <button onclick="logout()" class="text-sm text-gray-600 hover:text-gray-800">Logout</button>
    </div>

    <div class="flex h-screen">
        <!-- Desktop Sidebar -->
        <div class="desktop-sidebar w-16 bg-gray-800 flex flex-col items-center py-4 space-y-4">
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
            <div class="bg-white border-b px-6 py-4 sticky top-0 z-10">
                <div class="flex items-center space-x-4">
                    <!-- Global search -->
                    <div class="flex-1 max-w-md relative">
                        <input type="text"
                               placeholder="Search content..."
                               class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                               hx-get="/api/search/content"
                               hx-trigger="keyup changed delay:500ms"
                               hx-target="#search-results"
                               hx-vals='js:{"q": event.target.value}'
                               onfocus="showSearchResults()"
                               onblur="hideSearchResults()">
                        <div id="search-results" class="absolute top-full left-0 right-0 bg-white border rounded-lg shadow-lg z-20 hidden"></div>
                    </div>
                    <button onclick="logout()" class="px-4 py-2 text-gray-600 hover:text-gray-800 hidden md:block">
                        Logout
                    </button>
                </div>
            </div>

            <!-- Content area -->
            <div class="flex-1 flex flex-col md:flex-row overflow-hidden">
                <!-- Folders panel -->
                <div class="folder-panel w-full md:w-80 bg-white border-r flex flex-col">
                    <!-- Folders header -->
                    <div class="border-b px-4 py-3 flex items-center justify-between bg-gray-50">
                        <h2 class="font-semibold text-gray-800 flex items-center">
                            <svg class="w-5 h-5 mr-2 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"></path>
                            </svg>
                            Folders
                        </h2>
                        <button class="text-blue-600 hover:text-blue-800 px-3 py-1 text-sm font-medium rounded-md hover:bg-blue-50 transition-colors"
                                onclick="openCreateFolderModal()">
                            <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
                            </svg>
                            New
                        </button>
                    </div>

                    <!-- Breadcrumb -->
                    <div id="breadcrumb" class="px-4 py-3 border-b bg-gray-50">
                        <nav class="flex" aria-label="Breadcrumb">
                            <ol class="inline-flex items-center space-x-1 text-sm">
                                <li class="inline-flex items-center">
                                    <button class="text-gray-700 hover:text-blue-600 font-medium"
                                            onclick="navigateToRoot()">
                                        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"></path>
                                        </svg>
                                        Root
                                    </button>
                                </li>
                            </ol>
                        </nav>
                    </div>

                    <!-- Folder list -->
                    <div id="folder-list" class="flex-1 overflow-y-auto"
                         hx-get="/api/folders"
                         hx-trigger="load">
                        <div class="p-4 text-center text-gray-500">
                            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto mb-2"></div>
                            Loading folders...
                        </div>
                    </div>
                </div>

                <!-- Content panel -->
                <div class="content-panel flex-1 flex flex-col">
                    <!-- Content header -->
                    <div class="border-b px-6 py-4 flex items-center justify-between bg-gray-50">
                        <h2 id="content-title" class="font-semibold text-gray-800 flex items-center">
                            <svg class="w-5 h-5 mr-2 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                            </svg>
                            Select a folder to view contents
                        </h2>
                        <button id="new-content-btn" class="hidden text-blue-600 hover:text-blue-800 px-3 py-1 text-sm font-medium rounded-md hover:bg-blue-50 transition-colors"
                                onclick="openCreateContentModal()">
                            <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
                            </svg>
                            New Content
                        </button>
                    </div>

                    <!-- Content list -->
                    <div id="content-list" class="flex-1 overflow-y-auto p-6">
                        <div class="text-center text-gray-500 mt-12">
                            <svg class="w-16 h-16 mx-auto text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"></path>
                            </svg>
                            <h3 class="text-lg font-medium text-gray-900 mb-2">Welcome to your workspace</h3>
                            <p class="text-gray-500 mb-4">Select a folder from the left panel to view its contents</p>
                            <button onclick="openCreateFolderModal()" class="inline-flex items-center px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700 transition-colors">
                                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
                                </svg>
                                Create Your First Folder
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- Modals container -->
    <div id="modals"></div>

    <script>
        let currentFolderId = null;
        let currentFolderName = '';

        function selectFolder(folderId, folderName) {
            currentFolderId = folderId;
            currentFolderName = folderName;

            // Update content title
            document.getElementById('content-title').innerHTML =
                '<svg class="w-5 h-5 mr-2 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"></path></svg>' +
                folderName + ' Contents';

            document.getElementById('new-content-btn').classList.remove('hidden');

            // Add loading state
            document.getElementById('content-list').classList.add('loading');
            document.getElementById('folder-list').classList.add('loading');

            // Load folder contents
            htmx.ajax('GET', '/api/folders/' + folderId + '/content', '#content-list');

            // Load subfolders
            htmx.ajax('GET', '/api/folders?parent=' + folderId, '#folder-list');

            // Update breadcrumb
            htmx.ajax('GET', '/api/folders/' + folderId + '/path', '#breadcrumb');

            // Remove loading state after a short delay
            setTimeout(() => {
                document.getElementById('content-list').classList.remove('loading');
                document.getElementById('folder-list').classList.remove('loading');
            }, 300);
        }

        function navigateToRoot() {
            currentFolderId = null;
            currentFolderName = '';

            document.getElementById('content-title').innerHTML =
                '<svg class="w-5 h-5 mr-2 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>' +
                'Select a folder to view contents';

            document.getElementById('new-content-btn').classList.add('hidden');
            document.getElementById('content-list').innerHTML =
                '<div class="text-center text-gray-500 mt-12"><svg class="w-16 h-16 mx-auto text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"></path></svg><h3 class="text-lg font-medium text-gray-900 mb-2">Welcome to your workspace</h3><p class="text-gray-500 mb-4">Select a folder from the left panel to view its contents</p><button onclick="openCreateFolderModal()" class="inline-flex items-center px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700 transition-colors"><svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path></svg>Create Your First Folder</button></div>';

            // Reset breadcrumb
            document.getElementById('breadcrumb').innerHTML =
                '<nav class="flex" aria-label="Breadcrumb"><ol class="inline-flex items-center space-x-1 text-sm"><li class="inline-flex items-center"><button class="text-gray-700 hover:text-blue-600 font-medium" onclick="navigateToRoot()"><svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"></path></svg>Root</button></li></ol></nav>';

            // Load root folders
            htmx.ajax('GET', '/api/folders', '#folder-list');
        }

        function openCreateFolderModal() {
            htmx.ajax('GET', '/modals/create-folder?parent=' + (currentFolderId || ''), '#modals');
        }

        function openCreateContentModal() {
            if (currentFolderId) {
                htmx.ajax('GET', '/modals/create-content?folder=' + currentFolderId, '#modals');
            }
        }

        function editFolder(folderId) {
            htmx.ajax('GET', '/modals/edit-folder/' + folderId, '#modals');
        }

        function editContent(contentId) {
            htmx.ajax('GET', '/modals/edit-content/' + contentId, '#modals');
        }

        function showSearchResults() {
            const results = document.getElementById('search-results');
            if (results.innerHTML.trim()) {
                results.classList.remove('hidden');
            }
        }

        function hideSearchResults() {
            setTimeout(() => {
                document.getElementById('search-results').classList.add('hidden');
            }, 200);
        }

        function logout() {
            htmx.ajax('POST', '/logout', {target: 'body'});
        }

        function deleteFolder(folderId) {
            if (confirm('Are you sure you want to delete this folder? This will also delete all contents inside it.')) {
                htmx.ajax('DELETE', '/api/folders/' + folderId, {
                    target: '#folder-list',
                    swap: 'innerHTML',
                    callback: function() {
                        if (currentFolderId === folderId) {
                            navigateToRoot();
                        } else {
                            htmx.ajax('GET', '/api/folders' + (currentFolderId ? '?parent=' + currentFolderId : ''), '#folder-list');
                        }
                    }
                });
            }
        }

        function deleteContent(contentId) {
            if (confirm('Are you sure you want to delete this content?')) {
                htmx.ajax('DELETE', '/api/folders/' + currentFolderId + '/content/' + contentId, {
                    target: '#content-list',
                    swap: 'innerHTML',
                    callback: function() {
                        htmx.ajax('GET', '/api/folders/' + currentFolderId + '/content', '#content-list');
                    }
                });
            }
        }

        // Mobile menu toggle
        document.getElementById('mobile-menu-btn')?.addEventListener('click', function() {
            const folderPanel = document.querySelector('.folder-panel');
            folderPanel.classList.toggle('hidden');
        });

        // Add fade-in effect to dynamic content
        document.body.addEventListener('htmx:afterSwap', function(event) {
            if (event.target.id === 'folder-list' || event.target.id === 'content-list') {
                event.target.classList.add('fade-in');
            }
        });
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
			<div class="group p-4 border-b hover:bg-gray-50 transition-colors">
				<div class="flex items-center justify-between">
					<div class="flex items-center space-x-3 flex-1 cursor-pointer" onclick="selectFolder('%s', '%s')">
						<svg class="w-5 h-5 text-yellow-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"></path>
						</svg>
						<div class="flex-1 min-w-0">
							<h3 class="text-sm font-medium text-gray-900 truncate">%s</h3>
							<p class="text-xs text-gray-500 truncate">%s</p>
						</div>
					</div>
					<div class="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity">
						<button onclick="editFolder('%s')"
								class="p-1 text-gray-400 hover:text-blue-600 rounded"
								title="Edit folder">
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
							</svg>
						</button>
						<button onclick="deleteFolder('%s')"
								class="p-1 text-gray-400 hover:text-red-600 rounded"
								title="Delete folder">
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
							</svg>
						</button>
					</div>
				</div>
			</div>`,
			folder.ID, folder.Name, folder.Name, folder.Description, folder.ID, folder.ID)
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

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	parentIDStr := r.FormValue("parent_id")

	if name == "" {
		http.Error(w, "Folder name is required", http.StatusBadRequest)
		return
	}

	folder := models.Folder{
		Name:        name,
		Description: description,
		UserID:      userID,
	}

	// Set parent ID if provided
	if parentIDStr != "" {
		parentID, err := uuid.Parse(parentIDStr)
		if err != nil {
			http.Error(w, "Invalid parent ID", http.StatusBadRequest)
			return
		}
		folder.ParentID = &parentID
	}

	if err := h.folderService.CreateFolder(&folder); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return updated folder list
	var folders []models.Folder
	var err error
	if folder.ParentID != nil {
		folders, err = h.folderService.GetFoldersByParent(*folder.ParentID, userID)
	} else {
		folders, err = h.folderService.GetRootFolders(userID)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render folder list
	html := ""
	for _, f := range folders {
		html += fmt.Sprintf(`
			<div class="group p-4 border-b hover:bg-gray-50 transition-colors">
				<div class="flex items-center justify-between">
					<div class="flex items-center space-x-3 flex-1 cursor-pointer" onclick="selectFolder('%s', '%s')">
						<svg class="w-5 h-5 text-yellow-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"></path>
						</svg>
						<div class="flex-1 min-w-0">
							<h3 class="text-sm font-medium text-gray-900 truncate">%s</h3>
							<p class="text-xs text-gray-500 truncate">%s</p>
						</div>
					</div>
					<div class="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity">
						<button onclick="editFolder('%s')"
								class="p-1 text-gray-400 hover:text-blue-600 rounded"
								title="Edit folder">
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
							</svg>
						</button>
						<button onclick="deleteFolder('%s')"
								class="p-1 text-gray-400 hover:text-red-600 rounded"
								title="Delete folder">
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
							</svg>
						</button>
					</div>
				</div>
			</div>`,
			f.ID, f.Name, f.Name, f.Description, f.ID, f.ID)
	}

	if html == "" {
		html = `<div class="p-4 text-gray-500 text-center">No folders found</div>`
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}// UpdateFolder updates an existing folder
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
		contentTypeIcon := getContentTypeIcon(content.ContentType)
		html += fmt.Sprintf(`
			<div class="group bg-white p-4 rounded-lg border mb-3 hover:shadow-md transition-shadow">
				<div class="flex items-start justify-between">
					<div class="flex items-start space-x-3 flex-1">
						<div class="flex-shrink-0">
							%s
						</div>
						<div class="flex-1 min-w-0">
							<h3 class="font-medium text-gray-900 mb-1">%s</h3>
							<p class="text-sm text-gray-600 mb-2">%s</p>
							<div class="flex items-center space-x-4 text-xs text-gray-400">
								<span class="inline-flex items-center px-2 py-1 rounded-full bg-gray-100 text-gray-600">
									%s
								</span>
								<span>%s</span>
							</div>
						</div>
					</div>
					<div class="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity">
						<button onclick="editContent('%s')"
								class="p-2 text-gray-400 hover:text-blue-600 rounded-md"
								title="Edit content">
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
							</svg>
						</button>
						<button onclick="deleteContent('%s')"
								class="p-2 text-gray-400 hover:text-red-600 rounded-md"
								title="Delete content">
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
							</svg>
						</button>
					</div>
				</div>
			</div>`,
			contentTypeIcon, content.Title, content.Description, content.ContentType,
			content.CreatedAt.Format("Jan 2, 2006"), content.ID, content.ID)
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

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	body := r.FormValue("body")
	contentType := r.FormValue("content_type")

	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if contentType == "" {
		contentType = "text"
	}

	content := models.Content{
		Title:       title,
		Description: description,
		Body:        body,
		ContentType: contentType,
		FolderID:    folderID,
		UserID:      userID,
	}

	if err := h.contentService.CreateContent(&content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return updated content list
	contents, err := h.contentService.GetContentByFolderID(folderID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render content list
	html := ""
	for _, c := range contents {
		contentTypeIcon := getContentTypeIcon(c.ContentType)
		html += fmt.Sprintf(`
			<div class="group bg-white p-4 rounded-lg border mb-3 hover:shadow-md transition-shadow">
				<div class="flex items-start justify-between">
					<div class="flex items-start space-x-3 flex-1">
						<div class="flex-shrink-0">
							%s
						</div>
						<div class="flex-1 min-w-0">
							<h3 class="font-medium text-gray-900 mb-1">%s</h3>
							<p class="text-sm text-gray-600 mb-2">%s</p>
							<div class="flex items-center space-x-4 text-xs text-gray-400">
								<span class="inline-flex items-center px-2 py-1 rounded-full bg-gray-100 text-gray-600">
									%s
								</span>
								<span>%s</span>
							</div>
						</div>
					</div>
					<div class="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity">
						<button onclick="editContent('%s')"
								class="p-2 text-gray-400 hover:text-blue-600 rounded-md"
								title="Edit content">
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
							</svg>
						</button>
						<button onclick="deleteContent('%s')"
								class="p-2 text-gray-400 hover:text-red-600 rounded-md"
								title="Delete content">
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
							</svg>
						</button>
					</div>
				</div>
			</div>`,
			contentTypeIcon, c.Title, c.Description, c.ContentType,
			c.CreatedAt.Format("Jan 2, 2006"), c.ID, c.ID)
	}

	if html == "" {
		html = `<div class="text-gray-500 text-center">No content found in this folder</div>`
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
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

// Modal handlers
func (h *FolderHandler) CreateFolderModal(w http.ResponseWriter, r *http.Request) {
	parentID := r.URL.Query().Get("parent")

	html := `
<div id="modal-overlay" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" onclick="closeModal()">
	<div class="bg-white rounded-lg p-6 w-full max-w-md mx-4" onclick="event.stopPropagation()">
		<div class="flex justify-between items-center mb-4">
			<h2 class="text-lg font-semibold text-gray-900">Create New Folder</h2>
			<button onclick="closeModal()" class="text-gray-400 hover:text-gray-600">
				<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
				</svg>
			</button>
		</div>
		<form hx-post="/api/folders" hx-target="#folder-list" hx-swap="innerHTML" hx-on:htmx:after-request="closeModal(); refreshFolders()">
			<input type="hidden" name="parent_id" value="` + parentID + `">
			<div class="mb-4">
				<label for="name" class="block text-sm font-medium text-gray-700 mb-2">Folder Name</label>
				<input type="text" id="name" name="name" required
					   class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
					   placeholder="Enter folder name">
			</div>
			<div class="mb-6">
				<label for="description" class="block text-sm font-medium text-gray-700 mb-2">Description (Optional)</label>
				<textarea id="description" name="description" rows="3"
						  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
						  placeholder="Enter folder description"></textarea>
			</div>
			<div class="flex justify-end space-x-3">
				<button type="button" onclick="closeModal()"
						class="px-4 py-2 text-gray-600 border border-gray-300 rounded-md hover:bg-gray-50">
					Cancel
				</button>
				<button type="submit"
						class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">
					Create Folder
				</button>
			</div>
		</form>
	</div>
</div>

<script>
function closeModal() {
	document.getElementById('modal-overlay').remove();
}

function refreshFolders() {
	htmx.ajax('GET', '/api/folders' + (currentFolderId ? '?parent=' + currentFolderId : ''), '#folder-list');
}
</script>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (h *FolderHandler) EditFolderModal(w http.ResponseWriter, r *http.Request) {
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

	parentIDStr := ""
	if folder.ParentID != nil {
		parentIDStr = folder.ParentID.String()
	}

	html := fmt.Sprintf(`
<div id="modal-overlay" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" onclick="closeModal()">
	<div class="bg-white rounded-lg p-6 w-full max-w-md mx-4" onclick="event.stopPropagation()">
		<div class="flex justify-between items-center mb-4">
			<h2 class="text-lg font-semibold text-gray-900">Edit Folder</h2>
			<button onclick="closeModal()" class="text-gray-400 hover:text-gray-600">
				<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
				</svg>
			</button>
		</div>
		<form hx-put="/api/folders/%s" hx-target="#folder-list" hx-swap="innerHTML" hx-on:htmx:after-request="closeModal(); refreshFolders()">
			<input type="hidden" name="parent_id" value="%s">
			<div class="mb-4">
				<label for="name" class="block text-sm font-medium text-gray-700 mb-2">Folder Name</label>
				<input type="text" id="name" name="name" required value="%s"
					   class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
			</div>
			<div class="mb-6">
				<label for="description" class="block text-sm font-medium text-gray-700 mb-2">Description</label>
				<textarea id="description" name="description" rows="3"
						  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">%s</textarea>
			</div>
			<div class="flex justify-between">
				<button type="button" onclick="deleteFolder('%s')"
						class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700">
					Delete
				</button>
				<div class="space-x-3">
					<button type="button" onclick="closeModal()"
							class="px-4 py-2 text-gray-600 border border-gray-300 rounded-md hover:bg-gray-50">
						Cancel
					</button>
					<button type="submit"
							class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">
						Update
					</button>
				</div>
			</div>
		</form>
	</div>
</div>

<script>
function closeModal() {
	document.getElementById('modal-overlay').remove();
}

function refreshFolders() {
	htmx.ajax('GET', '/api/folders' + (currentFolderId ? '?parent=' + currentFolderId : ''), '#folder-list');
}

function deleteFolder(folderId) {
	if (confirm('Are you sure you want to delete this folder? This will also delete all contents inside it.')) {
		htmx.ajax('DELETE', '/api/folders/' + folderId, {target: '#folder-list', swap: 'innerHTML'});
		closeModal();
		refreshFolders();
	}
}
</script>`, folder.ID, parentIDStr, folder.Name, folder.Description, folder.ID)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (h *FolderHandler) CreateContentModal(w http.ResponseWriter, r *http.Request) {
	folderID := r.URL.Query().Get("folder")

	html := `
<div id="modal-overlay" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" onclick="closeModal()">
	<div class="bg-white rounded-lg p-6 w-full max-w-2xl mx-4" onclick="event.stopPropagation()">
		<div class="flex justify-between items-center mb-4">
			<h2 class="text-lg font-semibold text-gray-900">Create New Content</h2>
			<button onclick="closeModal()" class="text-gray-400 hover:text-gray-600">
				<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
				</svg>
			</button>
		</div>
		<form hx-post="/api/folders/` + folderID + `/content" hx-target="#content-list" hx-swap="innerHTML" hx-on:htmx:after-request="closeModal(); refreshContent()">
			<div class="mb-4">
				<label for="title" class="block text-sm font-medium text-gray-700 mb-2">Title</label>
				<input type="text" id="title" name="title" required
					   class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
					   placeholder="Enter content title">
			</div>
			<div class="mb-4">
				<label for="content_type" class="block text-sm font-medium text-gray-700 mb-2">Content Type</label>
				<select id="content_type" name="content_type"
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
					<option value="text">Text</option>
					<option value="note">Note</option>
					<option value="document">Document</option>
					<option value="link">Link</option>
				</select>
			</div>
			<div class="mb-4">
				<label for="description" class="block text-sm font-medium text-gray-700 mb-2">Description</label>
				<textarea id="description" name="description" rows="3"
						  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
						  placeholder="Enter content description"></textarea>
			</div>
			<div class="mb-6">
				<label for="body" class="block text-sm font-medium text-gray-700 mb-2">Content</label>
				<textarea id="body" name="body" rows="6"
						  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
						  placeholder="Enter the main content"></textarea>
			</div>
			<div class="flex justify-end space-x-3">
				<button type="button" onclick="closeModal()"
						class="px-4 py-2 text-gray-600 border border-gray-300 rounded-md hover:bg-gray-50">
					Cancel
				</button>
				<button type="submit"
						class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">
					Create Content
				</button>
			</div>
		</form>
	</div>
</div>

<script>
function closeModal() {
	document.getElementById('modal-overlay').remove();
}

function refreshContent() {
	if (currentFolderId) {
		htmx.ajax('GET', '/api/folders/' + currentFolderId + '/content', '#content-list');
	}
}
</script>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (h *FolderHandler) EditContentModal(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	contentID, err := uuid.Parse(chi.URLParam(r, "contentID"))
	if err != nil {
		http.Error(w, "Invalid content ID", http.StatusBadRequest)
		return
	}

	content, err := h.contentService.GetContentByID(contentID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	html := fmt.Sprintf(`
<div id="modal-overlay" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" onclick="closeModal()">
	<div class="bg-white rounded-lg p-6 w-full max-w-2xl mx-4" onclick="event.stopPropagation()">
		<div class="flex justify-between items-center mb-4">
			<h2 class="text-lg font-semibold text-gray-900">Edit Content</h2>
			<button onclick="closeModal()" class="text-gray-400 hover:text-gray-600">
				<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
				</svg>
			</button>
		</div>
		<form hx-put="/api/folders/%s/content/%s" hx-target="#content-list" hx-swap="innerHTML" hx-on:htmx:after-request="closeModal(); refreshContent()">
			<div class="mb-4">
				<label for="title" class="block text-sm font-medium text-gray-700 mb-2">Title</label>
				<input type="text" id="title" name="title" required value="%s"
					   class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
			</div>
			<div class="mb-4">
				<label for="content_type" class="block text-sm font-medium text-gray-700 mb-2">Content Type</label>
				<select id="content_type" name="content_type"
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
					<option value="text"%s>Text</option>
					<option value="note"%s>Note</option>
					<option value="document"%s>Document</option>
					<option value="link"%s>Link</option>
				</select>
			</div>
			<div class="mb-4">
				<label for="description" class="block text-sm font-medium text-gray-700 mb-2">Description</label>
				<textarea id="description" name="description" rows="3"
						  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">%s</textarea>
			</div>
			<div class="mb-6">
				<label for="body" class="block text-sm font-medium text-gray-700 mb-2">Content</label>
				<textarea id="body" name="body" rows="6"
						  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">%s</textarea>
			</div>
			<div class="flex justify-between">
				<button type="button" onclick="deleteContent('%s')"
						class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700">
					Delete
				</button>
				<div class="space-x-3">
					<button type="button" onclick="closeModal()"
							class="px-4 py-2 text-gray-600 border border-gray-300 rounded-md hover:bg-gray-50">
						Cancel
					</button>
					<button type="submit"
							class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">
						Update
					</button>
				</div>
			</div>
		</form>
	</div>
</div>

<script>
function closeModal() {
	document.getElementById('modal-overlay').remove();
}

function refreshContent() {
	if (currentFolderId) {
		htmx.ajax('GET', '/api/folders/' + currentFolderId + '/content', '#content-list');
	}
}

function deleteContent(contentId) {
	if (confirm('Are you sure you want to delete this content?')) {
		htmx.ajax('DELETE', '/api/folders/' + currentFolderId + '/content/' + contentId, {target: '#content-list', swap: 'innerHTML'});
		closeModal();
		refreshContent();
	}
}
</script>`,
		content.FolderID, content.ID, content.Title,
		getSelectedOption(content.ContentType, "text"),
		getSelectedOption(content.ContentType, "note"),
		getSelectedOption(content.ContentType, "document"),
		getSelectedOption(content.ContentType, "link"),
		content.Description, content.Body, content.ID)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func getSelectedOption(contentType, option string) string {
	if contentType == option {
		return " selected"
	}
	return ""
}

func getContentTypeIcon(contentType string) string {
	icons := map[string]string{
		"text":     `<svg class="w-5 h-5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>`,
		"note":     `<svg class="w-5 h-5 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>`,
		"document": `<svg class="w-5 h-5 text-purple-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>`,
		"link":     `<svg class="w-5 h-5 text-orange-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"></path></svg>`,
	}

	if icon, exists := icons[contentType]; exists {
		return icon
	}
	return icons["text"] // Default to text icon
}
