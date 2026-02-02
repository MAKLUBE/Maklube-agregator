package handlers

import "net/http"

func (h *Handler) adminHalalRequests(w http.ResponseWriter, r *http.Request) {
	
	// TODO: halal_verifications store + approve/reject
	w.Write([]byte("admin halal requests (TODO)"))
}
