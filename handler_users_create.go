package main

import (
	"encoding/json"
	"net/http"

	"github.com/Hessam-Emami/Chirpy/database"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handleUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if len(params.Email) == 0 {
		respondWithError(w, http.StatusBadRequest, "Must send an email!")
		return
	}

	if len(params.Password) == 0 {
		respondWithError(w, http.StatusBadRequest, "Must send a password!")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email, bcrypt.GenerateFromPassword(params.Password))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, database.User{
		ID:    user.ID,
		Email: user.Email,
	})

}
