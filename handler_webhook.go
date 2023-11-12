package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Hessam-Emami/Chirpy/internal/auth"
	"github.com/Hessam-Emami/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		}
	}

	apiKey, err := auth.GetApiKeyFromHeader(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find Api Key")
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Wrong api key")
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, struct{}{})
		return
	}

	_, err = cfg.DB.UpgradeChirpyRed(params.Data.UserID)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
