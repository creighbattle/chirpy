package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	type validResponse struct {
		Valid bool `json:"valid"`
	}
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(400)
		dat, _ := json.Marshal(errorResponse{Error: "something went wrong"})
		w.Write(dat)
		return
	}

	if len(params.Body) > 140 {
		w.WriteHeader(400)
		dat, _ := json.Marshal(errorResponse{Error: "Chirp is too long"})
		w.Write(dat)
		return
	}

	w.WriteHeader(200)
	dat, _ := json.Marshal(validResponse{Valid: true})
	w.Write(dat)
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