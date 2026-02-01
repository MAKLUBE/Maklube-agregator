package handlers

import "net/http"

func (h *Handler) partnerIncomingOrders(w http.ResponseWriter, r *http.Request) {
	// TODO: orders store и and getspisok
	
	w.Write([]byte("partner incoming orders"))
}
