package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits += 1
		next.ServeHTTP(w, r)
	})
}

func healthCheckHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	responseWriter.WriteHeader(200)
	responseWriter.Write([]byte("OK"))
}

func (cfg *apiConfig) numOfRequestsHandler(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits)))
}

func (cfg * apiConfig) resetRequestCountHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("Reset"))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err:= json.Marshal(payload)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(code)
    w.Write(response)
    return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error{
	return respondWithJSON(w, code, map[string]string{"error": msg})
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type validResponse struct {
		Cleaned_Body string `json:"cleaned_body"`
	}

	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	filteredText := replaceBadWords(params.Body)

	respondWithJSON(w, 200, validResponse{Cleaned_Body: filteredText})
}


func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))
	fileServer = http.StripPrefix("/app", fileServer)

	config := apiConfig{
		 fileserverHits: 0,
	}

	// Handle the main app route with the middleware
	mux.Handle("/app/", config.middlewareMetricsInc(fileServer))
	
	// Handle the metrics endpoint
	mux.HandleFunc("GET /admin/metrics", config.numOfRequestsHandler)

	// Handle the reset endpoint
	mux.HandleFunc("/api/reset", config.resetRequestCountHandler)
	
	// Handle health check
	mux.HandleFunc("GET /api/healthz", healthCheckHandler)

	// Validate Chirp
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	server := &http.Server{
		Addr: "localhost:8080",
		Handler: mux,
	}
	
	server.ListenAndServe()
}



func replaceBadWords(s string) string {
	strArr := strings.Split(s, " ")

	for i, word := range strArr {
		wordLowered := strings.ToLower(word)
		if wordLowered == "kerfuffle" || wordLowered == "sharbert" || wordLowered == "fornax" {
			strArr[i] = "****"
		}
	}

	return strings.Join(strArr, " ")
}