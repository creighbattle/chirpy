package main

import (
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {

	
	chirpId := r.PathValue("chirpID")
	chirpIdInt, _ := strconv.Atoi(chirpId)

	accessToken := r.Header.Get("Authorization")
	if len(accessToken) == 0 {
		respondWithError(w, http.StatusUnauthorized, "access token required")
		return
	}
	accessToken = accessToken[7:]

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
    	return cfg.jwtSecret, nil
	})

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	userId, err := token.Claims.GetSubject()
	userIdInt, _ := strconv.Atoi(userId)
	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Subject not found")
		return
	}

	err = cfg.DB.DeleteChirp(userIdInt, chirpIdInt)

	if err != nil {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})



}