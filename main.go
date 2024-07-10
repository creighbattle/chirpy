package main

import "net/http"

func myHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	responseWriter.WriteHeader(200)
	responseWriter.Write([]byte("OK"))
}


func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))
	fileServer = http.StripPrefix("/app", fileServer)

	mux.Handle("/app/*", fileServer)
	mux.Handle("/app/assets/logo.png", fileServer)

	

	mux.HandleFunc("/healthz", myHandler)

	server := &http.Server{
		Addr: "localhost:8080",
		Handler: mux,
	}


	server.ListenAndServe()
}