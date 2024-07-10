package main

import "net/http"


func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))

	mux.Handle("/", fileServer)
	mux.Handle("/assets/logo.png", fileServer)

	server := &http.Server{
		Addr: "localhost:8080",
		Handler: mux,
	}


	server.ListenAndServe()
}