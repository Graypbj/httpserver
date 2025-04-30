package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Graypbj/httpserver/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") == "" {
		respondWithError(w, http.StatusUnauthorized, "No token sent", errors.New("No token was sent"))
		return
	}

	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to find user ID", err)
		return
	}

	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to parse chirp id", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp was not found", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "User does not own chirp", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to delete chirp", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
