package handlers

import (
	"net/http"

	"mh-go/internal/services"

	"github.com/google/uuid"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := h.userService.GetByID(userID)
	if err != nil {
		writeError(w, "User not found", http.StatusNotFound)
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
    <!-- Navigation -->
    <nav class="bg-white shadow-sm border-b">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div class="flex justify-between h-16">
                <div class="flex items-center">
                    <h1 class="text-xl font-semibold text-gray-900">MH-Go</h1>
                </div>
                <div class="flex items-center space-x-4">
                    <span class="text-gray-700">` + user.Email + `</span>
                    <button hx-post="/logout"
                            class="bg-red-600 hover:bg-red-700 text-white px-3 py-2 rounded-md text-sm font-medium">
                        Logout
                    </button>
                </div>
            </div>
        </div>
    </nav>

    <div class="flex h-screen pt-16">
        <!-- Sidebar -->
        <div class="w-64 bg-white shadow-sm border-r">
            <div class="p-4">
                <div class="flex justify-between items-center mb-4">
                    <h2 class="text-lg font-medium text-gray-900">Folders</h2>
                    <button hx-get="/folders/new" hx-target="#modal" hx-swap="innerHTML"
                            class="bg-indigo-600 hover:bg-indigo-700 text-white px-2 py-1 rounded text-sm">
                        New
                    </button>
                </div>
                <div id="folders-list" hx-get="/folders" hx-trigger="load">
                    <!-- Folders will be loaded here -->
                </div>
            </div>
        </div>

        <!-- Main Content -->
        <div class="flex-1 flex flex-col">
            <div class="bg-white border-b p-4">
                <div class="flex justify-between items-center">
                    <h2 class="text-lg font-medium text-gray-900" id="content-title">Select a folder</h2>
                    <button hx-get="/content/new" hx-target="#modal" hx-swap="innerHTML"
                            class="bg-indigo-600 hover:bg-indigo-700 text-white px-3 py-2 rounded-md text-sm font-medium">
                        New Content
                    </button>
                </div>
            </div>
            <div class="flex-1 p-4 overflow-y-auto">
                <div id="content-list">
                    <div class="text-center text-gray-500 mt-8">
                        <p>Select a folder to view content</p>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- Modal -->
    <div id="modal" class="hidden fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
        <!-- Modal content will be loaded here -->
    </div>

    <script>
        // Show modal when content is loaded
        document.body.addEventListener('htmx:afterSwap', function(event) {
            if (event.target.id === 'modal' && event.target.innerHTML.trim()) {
                event.target.classList.remove('hidden');
            }
        });

        // Hide modal
        function hideModal() {
            document.getElementById('modal').classList.add('hidden');
            document.getElementById('modal').innerHTML = '';
        }

        // Close modal on outside click
        document.getElementById('modal').addEventListener('click', function(event) {
            if (event.target === this) {
                hideModal();
            }
        });
    </script>
</body>
</html>`
	w.Write([]byte(html))
}


