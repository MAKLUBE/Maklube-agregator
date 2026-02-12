package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"
)

func (h *Handler) verifyForm(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("email")))
	h.render(w, r, "verify.tmpl", &templateData{
		Form: map[string]string{"email": email},
		User: h.currentUser(r),
	})
}

func (h *Handler) verifyPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(strings.ToLower(r.PostForm.Get("email")))
	code := strings.TrimSpace(r.PostForm.Get("code"))

	if email == "" || code == "" {
		h.render(w, r, "verify.tmpl", &templateData{
			Form: map[string]string{"error": "Fill all fields", "email": email},
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err := h.App.Users.VerifyByEmailAndCode(ctx, email, code, time.Now().UTC())
	if err != nil {
		h.render(w, r, "verify.tmpl", &templateData{
			Form: map[string]string{"error": "Invalid or expired code", "email": email},
		})
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
