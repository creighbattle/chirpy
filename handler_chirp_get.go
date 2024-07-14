package main

import (
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handlerChripRetrieve(w http.ResponseWriter, r *http.Request) {

	dbChrips, err := cfg.DB.LoadDB()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Coudn't retrieve chirps")
		return
	}
	pathValue := r.PathValue("chirpID")


	pathValueInt, _ := strconv.Atoi(pathValue)

	val, ok := dbChrips.Chirps[pathValueInt]

	

	if !ok {
		respondWithError(w, http.StatusNotFound, "The Chirp does not exist")
		return
	}

	respondWithJSON(w, http.StatusOK, val)




}