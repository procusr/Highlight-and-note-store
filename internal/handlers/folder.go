package handlers

import (
	"encoding/json"
	"net/http"

	"mh-go/internal/models"
	"mh-go/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type FolderHandler struct {
	folderService *services.FolderService
}

func NewFolderHandler(folderService *services.FolderService) *FolderHandler {
	return &FolderHandler{folderService: folderService}
}

func (h *FolderHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	folders, err := h.folderService.GetByUserID(userID)
	if err != nil {
		writeError(w, "Failed to get folders", http.StatusInternalServerError)
		return
	}

	html := `<div class="space-y-2">`
	for _, folder := range folders {
		html += `
		<div class="p-3 hover:bg-gray-50 rounded-md cursor-pointer border-l-4 border-transparent hover:border-indigo-500"
		     hx-get="/content?folder_id=` + folder.ID.String() + `" hx-target="#content-list" hx-swap="innerHTML"
		     onclick="document.getElementById('content-title').innerText = '` + folder.Name + `'">
			<div class="flex justify-between items-center">
				<span class="font-medium text-gray-900">` + folder.Name + `</span>
				` + func() string {
			if folder.IsDefault {
				return `<span class="text-xs text-gray-500">Default</span>`
			}
			return `<button hx-delete="/folders/` + folder.ID.String() + `" hx-confirm="Are you sure?"
					class="text-red-600 hover:text-red-800 text-xs">Delete</button>`
		}() + `
			</div>
			` + func() string {
			if folder.Description != nil {
				return `<p class="text-sm text-gray-600 mt-1">` + *folder.Description + `</p>`
			}
			return ""
		}() + `
		</div>`
	}
	html += `</div>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (h *FolderHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Try form data
		req.Name = r.FormValue("name")
		description := r.FormValue("description")
		if description != "" {
			req.Description = &description
		}
	}

	if req.Name == "" {
		writeError(w, "Name is required", http.StatusBadRequest)
		return
	}

	folder, err := h.folderService.Create(userID, &req)
	if err != nil {
		writeError(w, "Failed to create folder", http.StatusInternalServerError)
		return
	}

	writeJSON(w, folder)
}

func (h *FolderHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	folderID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	folder, err := h.folderService.GetByID(folderID, userID)
	if err != nil {
		writeError(w, "Folder not found", http.StatusNotFound)
		return
	}

	writeJSON(w, folder)
}

func (h *FolderHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	folderID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	folder, err := h.folderService.Update(folderID, userID, &req)
	if err != nil {
		writeError(w, "Failed to update folder", http.StatusInternalServerError)
		return
	}

	writeJSON(w, folder)
}

func (h *FolderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	folderID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	err = h.folderService.Delete(folderID, userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
