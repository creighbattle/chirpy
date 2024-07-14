package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("Authorization")
	refreshToken = refreshToken[7:]

	dbStructure, err  := cfg.DB.LoadDB()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	users := dbStructure.Users
	found := false
	id := 0

	for _, user := range users {
		
		if user.RefreshToken == refreshToken {
			
			found = true
			parsedTime, err := time.Parse(time.RFC3339, user.Exp)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
			if time.Now().After(parsedTime) {
				respondWithError(w, http.StatusUnauthorized, "Refresh token has expired")
				return
			}
			id = user.ID
			break
		}
	}

	if !found {
		respondWithError(w, http.StatusUnauthorized, "Refresh token does not exist")
		return
	}

	currentTime := time.Now().UTC()

	jwtRegisteredClaims := jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(currentTime),
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(3600) * time.Second)),
		Subject: strconv.Itoa(id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtRegisteredClaims)

	signedJwtToken, err := token.SignedString(cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not sign token")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{Token string `json:"token"`}{
		Token: signedJwtToken,
	})

}