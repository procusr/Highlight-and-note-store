package handlers

import (
	"encoding/json"
	"net/http"

	"mh-go/internal/models"
	"mh-go/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ContentHandler struct {
	contentService *services.ContentService
}

func NewContentHandler(contentService *services.ContentService) *ContentHandler {
	return &ContentHandler{contentService: contentService}
}

func (h *ContentHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var folderID *uuid.UUID
	if folderIDStr := r.URL.Query().Get("folder_id"); folderIDStr != "" {
		if id, err := uuid.Parse(folderIDStr); err == nil {
			folderID = &id
		}
	}

	contents, err := h.contentService.GetByUserID(userID, folderID)
	if err != nil {
		writeError(w, "Failed to get content", http.StatusInternalServerError)
		return
	}

	html := `<div class="space-y-4">`
	if len(contents) == 0 {
		html += `<div class="text-center text-gray-500 mt-8">
			<p>No content in this folder</p>
		</div>`
	} else {
		for _, content := range contents {
			html += `
			<div class="bg-white p-4 rounded-lg shadow-sm border hover:shadow-md transition-shadow">
				<div class="flex justify-between items-start mb-2">
					<div class="flex-1">
						` + func() string {
				if content.Title != nil {
					return `<h3 class="text-lg font-medium text-gray-900 mb-2">` + *content.Title + `</h3>`
				}
				return ""
			}() + `
						<div class="text-gray-700 mb-3">
							` + func() string {
				if len(content.Content) > 200 {
					return content.Content[:200] + "..."
				}
				return content.Content
			}() + `
						</div>
						<div class="flex items-center space-x-4 text-sm text-gray-500">
							<span>` + content.CreatedAt.Format("Jan 2, 2006") + `</span>
							` + func() string {
				if content.IsFavorite {
					return `<span class="text-yellow-500">★ Favorite</span>`
				}
				return ""
			}() + `
						</div>
					</div>
					<div class="flex space-x-2 ml-4">
						<button hx-post="/content/` + content.ID.String() + `/favorite" hx-target="this" hx-swap="outerHTML"
								class="text-gray-400 hover:text-yellow-500">
							` + func() string {
				if content.IsFavorite {
					return "★"
				}
				return "☆"
			}() + `
						</button>
						<button hx-get="/content/` + content.ID.String() + `/edit" hx-target="#modal" hx-swap="innerHTML"
								class="text-indigo-600 hover:text-indigo-800">
							Edit
						</button>
						<button hx-delete="/content/` + content.ID.String() + `" hx-confirm="Are you sure?"
								class="text-red-600 hover:text-red-800">
							Delete
						</button>
					</div>
				</div>
			</div>`
		}
	}
	html += `</div>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (h *ContentHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Try form data
		req.Content = r.FormValue("content")
		title := r.FormValue("title")
		if title != "" {
			req.Title = &title
		}
		anchor := r.FormValue("anchor")
		if anchor != "" {
			req.Anchor = &anchor
		}
		if folderIDStr := r.FormValue("folder_id"); folderIDStr != "" {
			if id, err := uuid.Parse(folderIDStr); err == nil {
				req.FolderID = &id
			}
		}
	}

	if req.Content == "" {
		writeError(w, "Content is required", http.StatusBadRequest)
		return
	}

	content, err := h.contentService.Create(userID, &req)
	if err != nil {
		writeError(w, "Failed to create content", http.StatusInternalServerError)
		return
	}

	writeJSON(w, content)
}

func (h *ContentHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	contentID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, "Invalid content ID", http.StatusBadRequest)
		return
	}

	content, err := h.contentService.GetByID(contentID, userID)
	if err != nil {
		writeError(w, "Content not found", http.StatusNotFound)
		return
	}

	writeJSON(w, content)
}

func (h *ContentHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	contentID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, "Invalid content ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	content, err := h.contentService.Update(contentID, userID, &req)
	if err != nil {
		writeError(w, "Failed to update content", http.StatusInternalServerError)
		return
	}

	writeJSON(w, content)
}

func (h *ContentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	contentID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, "Invalid content ID", http.StatusBadRequest)
		return
	}

	err = h.contentService.Delete(contentID, userID)
	if err != nil {
		writeError(w, "Failed to delete content", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ContentHandler) Archive(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	contentID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, "Invalid content ID", http.StatusBadRequest)
		return
	}

	err = h.contentService.Archive(contentID, userID)
	if err != nil {
		writeError(w, "Failed to archive content", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ContentHandler) ToggleFavorite(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == uuid.Nil {
		writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	contentID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, "Invalid content ID", http.StatusBadRequest)
		return
	}

	err = h.contentService.ToggleFavorite(contentID, userID)
	if err != nil {
		writeError(w, "Failed to toggle favorite", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
