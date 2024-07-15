package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {

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