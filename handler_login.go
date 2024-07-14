package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerLogin (w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
	}

	type response struct {
		Email string `json:"email"`
		ID int `json:"id"`
		Token string `json:"token"`
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
		fmt.Println(user.Password)
		respondWithError(w, http.StatusUnauthorized, "Passwords do not match")
		return
	}

	// Get the current UTC time
	currentTime := time.Now().UTC()

	expireTime := params.ExpiresInSeconds

	if expireTime == 0 || expireTime > 86400 {
		expireTime = 86400
	}

	
	jwtRegisteredClaims := jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(currentTime),
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(expireTime) * time.Second)),
		Subject: strconv.Itoa(id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtRegisteredClaims)

	signedJwtToken, err := token.SignedString(cfg.jwtSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not sign token")
		return
	}

	
	respondWithJSON(w, http.StatusOK, response{Email: user.Email, ID: id, Token: signedJwtToken})

	

}