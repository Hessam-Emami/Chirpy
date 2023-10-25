package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Hessam-Emami/Chirpy/database"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type User struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

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
		respondWithError(w, http.StatusBadRequest, "Must send an password!")
		return
	}
	users, err := cfg.DB.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	user := database.User{}
	for _, v := range users {
		if v.Email == params.Email {
			user = v
		}
	}
	log.Printf("Couldn't find")
	if len(user.Email) == 0 {
		respondWithError(w, http.StatusNotFound, "Email doesn't match")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Password or email doesn't match")
	}
	respondWithJSON(w, http.StatusOK, User{
		ID:    user.ID,
		Email: user.Email,
	})

}
