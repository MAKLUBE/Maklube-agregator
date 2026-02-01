package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/MAKLUBE/AP1_Final_Project/internal/app"
	"github.com/MAKLUBE/AP1_Final_Project/internal/models"
	"github.com/justinas/alice"
)

type Middleware struct {
	App *app.Application
}

func New(a *app.Application) *Middleware {
	return &Middleware{App: a}
}

func (m *Middleware) Chain(next http.Handler) http.Handler {
	return alice.New(
		m.recoverPanic,
		m.secureHeaders,
		m.logRequest,
		m.authenticate,
	).Then(next)
}

func (m *Middleware) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.App.Logger.Printf("panic: %v\n%s", err, debug.Stack())
				http.Error(w, "Server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		m.App.Logger.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

//Auth helpers

func (m *Middleware) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		if err != nil || c.Value == "" {
			next.ServeHTTP(w, r)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		sess, err := m.App.Sessions.FindByToken(ctx, c.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		u, err := m.App.Users.FindByID(ctx, sess.UserID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		r = r.WithContext(app.ContextSetUser(r.Context(), u))
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uAny, ok := app.ContextGetUser(r.Context())
		if !ok || uAny == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RequireRole(role models.UserRole, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uAny, ok := app.ContextGetUser(r.Context())
		if !ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		u, ok := uAny.(*models.User)
		if !ok || u.Role != role {
			http.Error(w, fmt.Sprintf("Forbidden (need role %s)", role), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
