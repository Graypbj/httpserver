package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Graypbj/httpserver/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsersRed(w http.ResponseWriter, r *http.Request) {
	type UserID struct {
		UserID uuid.UUID `json:"user_id"`
	}

	type parameters struct {
		Event string `json:"event"`
		Data  UserID `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error getting API Key", err)
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "API Key does not match", errors.New("API Key does not match"))
		return
	}

	var req parameters
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	if req.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := cfg.db.UpgradeToRed(r.Context(), req.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error upgrading to red", err)
		return
	}

	if userID != req.Data.UserID {
		respondWithError(w, http.StatusNotFound, "UserID not found", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
