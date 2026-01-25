package routes

import (
	"net/http"

	"AP1_Final_Project/internal/handlers"
)

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/status", handlers.Status)
	mux.HandleFunc("/problems", handlers.ListProblems)
	mux.HandleFunc("/submit", handlers.SubmitSolution)

	return mux
}
