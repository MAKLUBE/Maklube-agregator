package handlers

import "net/http"

func Status(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("GoContester backend is running"))
}
