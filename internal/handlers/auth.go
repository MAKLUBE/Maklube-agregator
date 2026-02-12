package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/MAKLUBE/AP1_Final_Project/internal/app"
	"github.com/MAKLUBE/AP1_Final_Project/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) registerForm(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "register.tmpl", &templateData{User: h.currentUser(r)})
}

func (h *Handler) registerPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(r.PostForm.Get("username"))
	email := strings.TrimSpace(strings.ToLower(r.PostForm.Get("email")))
	password := r.PostForm.Get("password")
	role := strings.TrimSpace(r.PostForm.Get("role"))

	if username == "" || email == "" || password == "" {
		h.render(w, r, "register.tmpl", &templateData{
			Form: map[string]string{"error": "Fill all fields"},
		})
		return
	}

	userRole := models.RoleCustomer
	if role == "partner" {
		userRole = models.RolePartner
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		h.serverError(w, err)
		return
	}

	u := &models.User{
		Role:         userRole,
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    time.Now().UTC(),
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if _, err := h.App.Users.FindByEmail(ctx, email); err == nil {
		h.render(w, r, "register.tmpl", &templateData{
			Form: map[string]string{"error": "Email already used"},
		})
		return
	}

	if err := h.App.Users.Insert(ctx, u); err != nil {
		h.serverError(w, err)
		return
	}

	code, err := new6DigitCode()
	if err != nil {
		h.serverError(w, err)
		return
	}

	expires := time.Now().UTC().Add(10 * time.Minute)

	if err := h.App.Users.SetVerification(ctx, u.ID, code, expires); err != nil {
		h.serverError(w, err)
		return
	}

	subject := "Makluber email verification code"
	body := "Your verification code is: " + code + "\n\nIt expires in 10 minutes."

	if err := h.App.Mailer.Send(u.Email, subject, body); err != nil {
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/verify?email="+url.QueryEscape(u.Email), http.StatusSeeOther)
}

func (h *Handler) loginForm(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "login.tmpl", &templateData{User: h.currentUser(r)})
}

func (h *Handler) loginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(strings.ToLower(r.PostForm.Get("email")))
	password := r.PostForm.Get("password")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	u, err := h.App.Users.FindByEmail(ctx, email)
	if err != nil {
		h.render(w, r, "login.tmpl", &templateData{
			Form: map[string]string{"error": "Invalid credentials"},
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)); err != nil {
		h.render(w, r, "login.tmpl", &templateData{
			Form: map[string]string{"error": "Invalid credentials"},
		})
		return
	}

	if !u.Verified {
		h.render(w, r, "login.tmpl", &templateData{
			Form: map[string]string{"error": "Please verify your email before login"},
		})
		return
	}

	token, err := newToken(32)
	if err != nil {
		h.serverError(w, err)
		return
	}

	sessDuration := 7 * 24 * time.Hour

	if u.Role == models.RolePartner {
		sessDuration = 30 * 24 * time.Hour
	} else if u.Role == models.RoleAdmin {
		sessDuration = 1 * time.Hour
	}

	sess := &models.Session{
		Token:     token,
		UserID:    u.ID,
		ExpiresAt: time.Now().UTC().Add(sessDuration),
	}
	if err := h.App.Sessions.Create(ctx, sess); err != nil {
		h.serverError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) logoutPost(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err == nil && c.Value != "" {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		_ = h.App.Sessions.Delete(ctx, c.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func newToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (h *Handler) currentUser(r *http.Request) *models.User {
	uAny, ok := app.ContextGetUser(r.Context())
	if !ok || uAny == nil {
		return nil
	}
	u, _ := uAny.(*models.User)
	return u
}

func new6DigitCode() (string, error) {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	n := int(b[0])<<16 | int(b[1])<<8 | int(b[2])
	code := n % 1000000
	return fmt.Sprintf("%06d", code), nil
}
