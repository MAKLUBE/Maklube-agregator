package handlers

import (
	"net/http"

	"github.com/MAKLUBE/AP1_Final_Project/internal/models"
)

type templateData struct {
	User *models.User
	Data any
	Form map[string]string
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "home.tmpl", &templateData{
		User: h.currentUser(r),
	})
}
