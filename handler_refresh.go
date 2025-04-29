package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Graypbj/httpserver/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type resp struct {
		Token string `json:"token"`
	}

	if r.Header.Get("Authorization") == "" {
		respondWithError(w, http.StatusBadRequest, "Missing refresh token", errors.New("No refresh token was sent"))
		return
	}
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	dbToken, err := cfg.db.GetToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error finding refresh token in database", err)
		return
	}

	if dbToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token has been revoked", errors.New("Refresh token has been revoked"))
		return
	}

	if time.Now().After(dbToken.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token has expired", errors.New("Refresh token has expired"))
		return
	}

	user, err := cfg.db.GetUserByRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving user", err)
		return
	}

	newAccessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error created access token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, resp{
		Token: newAccessToken,
	})
}
