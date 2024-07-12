package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerLogin (w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		Email string `json:"email"`
		ID int `json:"id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	db, _ := cfg.DB.LoadDB()

	id, ok := db.Emails[params.Email]
	if !ok {
		respondWithError(w, http.StatusNotFound, "Email does not exist")
		return
	}

	user, ok := db.Users[id]
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "Error getting user")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Passwords do not match")
		return
	}

	respondWithJSON(w, http.StatusOK, response{Email: user.Email, ID: id})

	

}