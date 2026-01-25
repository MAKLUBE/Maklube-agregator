package main

import (
	"AP1_Final_Project/internal/routes"
	"log"
	"net/http"
)

func main() {
	mux := routes.RegisterRoutes()

	server := &http.Server{
		Addr:    ":4000",
		Handler: mux,
	}

	log.Println("Starting server on :4000")
	err := server.ListenAndServe()
	log.Fatal(err)
}
