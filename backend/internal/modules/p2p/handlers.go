package p2p

import (
	"encoding/json"
	"net/http"
)

// Handlers holds HTTP handlers for the P2P module.
type Handlers struct {
	mgr *P2PManager
}

// NewHandlers creates new P2P API handlers.
func NewHandlers(mgr *P2PManager) *Handlers {
	return &Handlers{mgr: mgr}
}

// RegisterRoutes registers all P2P routes.
func (h *Handlers) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/p2p/status", h.statusHandler)
	mux.HandleFunc("GET /api/p2p/keys", h.keysHandler)
	mux.HandleFunc("GET /api/p2p/contacts", h.contactsHandler)
	mux.HandleFunc("POST /api/p2p/contacts", h.addContactHandler)
	mux.HandleFunc("DELETE /api/p2p/contacts/{id}", h.deleteContactHandler)
	mux.HandleFunc("GET /api/p2p/contacts/{id}/messages", h.messagesHandler)
	mux.HandleFunc("POST /api/p2p/contacts/{id}/send", h.sendHandler)
}

func (h *Handlers) statusHandler(w http.ResponseWriter, r *http.Request) {
	p2pWriteJSON(w, http.StatusOK, h.mgr.GetStatus())
}

func (h *Handlers) keysHandler(w http.ResponseWriter, r *http.Request) {
	p2pWriteJSON(w, http.StatusOK, h.mgr.GetKeyPair())
}

func (h *Handlers) contactsHandler(w http.ResponseWriter, r *http.Request) {
	contacts := h.mgr.GetContacts()
	p2pWriteJSON(w, http.StatusOK, map[string]interface{}{
		"contacts": contacts,
		"count":    len(contacts),
	})
}

func (h *Handlers) addContactHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Alias     string `json:"alias"`
		PublicKey string `json:"public_key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		p2pWriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	contact, err := h.mgr.AddContact(body.Alias, body.PublicKey)
	if err != nil {
		p2pWriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	p2pWriteJSON(w, http.StatusCreated, contact)
}

func (h *Handlers) deleteContactHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.mgr.DeleteContact(id); err != nil {
		p2pWriteError(w, http.StatusNotFound, err.Error())
		return
	}
	p2pWriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handlers) messagesHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	msgs := h.mgr.GetMessages(id)
	p2pWriteJSON(w, http.StatusOK, map[string]interface{}{
		"contact_id": id,
		"messages":   msgs,
		"count":      len(msgs),
	})
}

func (h *Handlers) sendHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var body struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		p2pWriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	msg, err := h.mgr.SendMessage(id, body.Content)
	if err != nil {
		p2pWriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	p2pWriteJSON(w, http.StatusCreated, msg)
}

func p2pWriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func p2pWriteError(w http.ResponseWriter, status int, message string) {
	p2pWriteJSON(w, status, map[string]interface{}{"error": message, "status": status})
}
