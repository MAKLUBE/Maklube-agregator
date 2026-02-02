package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handler) serverError(w http.ResponseWriter, err error) {
	h.App.Logger.Println(err)
	http.Error(w, "Server error", http.StatusInternalServerError)
}

func (h *Handler) clientError(w http.ResponseWriter, status int) {
	http.Error(w, fmt.Sprintf("Client error: %d", status), status)
}

func (h *Handler) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := h.App.Templates[name]
	if !ok {
		h.serverError(w, fmt.Errorf("template %s not found", name))
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := ts.ExecuteTemplate(w, "base", td); err != nil {
		h.serverError(w, err)
	}
}
