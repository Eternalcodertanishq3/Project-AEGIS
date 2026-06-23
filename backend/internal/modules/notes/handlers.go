package notes

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
)

// Handlers holds the HTTP handler functions for the notes module.
type Handlers struct {
	notes *NotesManager
}

// NewHandlers creates new notes API handlers.
func NewHandlers(notes *NotesManager) *Handlers {
	return &Handlers{notes: notes}
}

// RegisterRoutes registers all notes module routes.
func (h *Handlers) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/notes", h.listHandler)
	mux.HandleFunc("GET /api/notes/{id}", h.getHandler)
	mux.HandleFunc("POST /api/notes", h.createHandler)
	mux.HandleFunc("PUT /api/notes/{id}", h.updateHandler)
	mux.HandleFunc("DELETE /api/notes/{id}", h.deleteHandler)
	mux.HandleFunc("POST /api/notes/{id}/pin", h.pinHandler)
}

func (h *Handlers) listHandler(w http.ResponseWriter, r *http.Request) {
	notes, err := h.notes.List()
	if err != nil {
		notesWriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	notesWriteJSON(w, http.StatusOK, map[string]interface{}{
		"notes": notes,
		"count": len(notes),
	})
}

func (h *Handlers) getHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	note, err := h.notes.Get(id)
	if err != nil {
		notesWriteError(w, http.StatusNotFound, err.Error())
		return
	}
	notesWriteJSON(w, http.StatusOK, note)
}

type createNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Tags    string `json:"tags"`
}

func (h *Handlers) createHandler(w http.ResponseWriter, r *http.Request) {
	var req createNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		notesWriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	id := generateID()
	note, err := h.notes.Create(id, req.Title, req.Content, req.Tags)
	if err != nil {
		notesWriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	notesWriteJSON(w, http.StatusCreated, note)
}

func (h *Handlers) updateHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		notesWriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	note, err := h.notes.Update(id, req.Title, req.Content, req.Tags)
	if err != nil {
		notesWriteError(w, http.StatusNotFound, err.Error())
		return
	}
	notesWriteJSON(w, http.StatusOK, note)
}

func (h *Handlers) deleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.notes.Delete(id); err != nil {
		notesWriteError(w, http.StatusNotFound, err.Error())
		return
	}
	notesWriteJSON(w, http.StatusOK, map[string]string{"message": "note deleted"})
}

func (h *Handlers) pinHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	note, err := h.notes.TogglePin(id)
	if err != nil {
		notesWriteError(w, http.StatusNotFound, err.Error())
		return
	}
	notesWriteJSON(w, http.StatusOK, note)
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func notesWriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func notesWriteError(w http.ResponseWriter, status int, message string) {
	notesWriteJSON(w, status, map[string]interface{}{"error": message, "status": status})
}
