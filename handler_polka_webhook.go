package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {

	apiKey := r.Header.Get("Authorization")
	if len(apiKey) < 8 {
		respondWithError(w, http.StatusUnauthorized, "api key required")
		return
	}
	apiKey = apiKey[7:]

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "invalid api key")
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	userId := params.Data.UserID
	event := params.Event

	if event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, struct{}{})
		return
	}

	err = cfg.DB.UpdateUserSubscription(userId)

	if err != nil && err.Error() == "no user" {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusNoContent, struct{}{})

}