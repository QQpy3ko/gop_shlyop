package api

import (
	"encoding/json"
	"errors"
	"gop_shlyop/internal/errs"
	"github.com/rs/zerolog/log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, r, errs.NewErrBadRequest("Only GET method is allowed"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
}

func respondWithError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Error processing request %s %s: %v", r.Method, r.URL.Path, err)

	var statusCode int
	var errMsg string

	var badRequestErr *errs.ErrBadRequest
	var notFoundErr *errs.ErrNotFound

	switch {
	case errors.As(err, &badRequestErr):
		statusCode = http.StatusBadRequest
		errMsg = badRequestErr.Error()
	case errors.As(err, &notFoundErr):
		statusCode = http.StatusNotFound
		errMsg = notFoundErr.Error()
	default:
		statusCode = http.StatusInternalServerError
		errMsg = "internal server error"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: errMsg})
}

func ExampleErrorHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, r, errs.NewErrNotFound("requested resource was not found"))
}
 