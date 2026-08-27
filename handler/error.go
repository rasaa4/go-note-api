package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"note-api/repository"
	"note-api/service"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error: err.Error(),
	})
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidID),
		errors.Is(err, service.ErrTitleRequired),
		errors.Is(err, service.ErrContentRequired),
		errors.Is(err, service.ErrTitleTooLong),
		errors.Is(err, service.ErrNoteRequired):

		writeError(w, http.StatusBadRequest, err)

	case errors.Is(err, repository.ErrNoteNotFound):
		writeError(w, http.StatusNotFound, err)

	default:
		writeError(
			w,
			http.StatusInternalServerError,
			errors.New("internal server error"),
		)
	}
}
