package main

import (
	"net/http"
	"sort"

	"github.com/Graypbj/httpserver/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	var dbChirps []database.Chirp
	var err error
	s := r.URL.Query().Get("author_id")
	sortVal := r.URL.Query().Get("sort")
	if s != "" {
		userID, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error parsing author_id", err)
			return
		}
		dbChirps, err = cfg.db.GetChirpsByUserID(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		}
	} else {
		dbChirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
			return
		}
	}

	chirps := []Chirp{}

	for _, dbChirp := range dbChirps {
		chirp := Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
		chirps = append(chirps, chirp)
	}

	if sortVal == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
