package main

import "net/http"

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("Authorization")
	refreshToken = refreshToken[7:]

	err := cfg.DB.RevokeRefreshToken(refreshToken)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}