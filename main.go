package main

import "net/http"

func main() {
	serveMux := http.NewServeMux()
	newServer := http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	newServer.ListenAndServe()
}
