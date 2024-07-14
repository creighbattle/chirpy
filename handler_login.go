package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
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
		RefreshToken string `json:"refresh_token"`
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

	// Get the current UTC time
	currentTime := time.Now().UTC()

	expireTime := params.ExpiresInSeconds

	if expireTime == 0 || expireTime > 3600 {
		expireTime = 3600
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

	c := 32
	b := make([]byte, c)
	_, err = rand.Read(b)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	futureTime := currentTime.Add(time.Duration(1440) * time.Hour)
	futureTimeString := futureTime.Format(time.RFC3339)

	err = cfg.DB.UpdateRefreshToken(id, hex.EncodeToString(b), futureTimeString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	
	respondWithJSON(w, http.StatusOK, response{Email: user.Email, ID: id, Token: signedJwtToken, RefreshToken: hex.EncodeToString(b)})

	

}