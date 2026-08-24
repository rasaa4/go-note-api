package handler

import (
	"encoding/json"
	"net/http"
	"note-api/service"
	"strconv"
)

type NoteHandler struct {
	service *service.NoteService
}

func NewNoteHandler(service *service.NoteService) *NoteHandler {
	return &NoteHandler{service: service}
}

type CreateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
type GetByIDRequest struct {
	ID int `json:"id"`
	CreateNoteRequest
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateNoteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	note, err := h.service.CreateNote(ctx, req.Title, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)

}
func (h *NoteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	note, err := h.service.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(note)

}
func (h *NoteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	notes, err := h.service.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notes)

}
