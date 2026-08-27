package handler

import (
	"errors"
	"net/http"
	"note-api/repository"
	"note-api/service"
)

func handleError(w http.ResponseWriter, err error) {

	if errors.Is(err, service.ErrInvalidID) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if errors.Is(err, service.ErrTitleRequired) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if errors.Is(err, service.ErrContentRequired) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if errors.Is(err, service.ErrTitleTooLong) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if errors.Is(err, service.ErrNoteRequired) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if errors.Is(err, repository.ErrNoteNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	http.Error(w, "internal server error", http.StatusInternalServerError)
}
