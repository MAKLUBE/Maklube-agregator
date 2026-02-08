package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) adminHalalRequests(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	items, err := h.App.Halal.ListPending(ctx)
	if err != nil {
		h.serverError(w, err)
		return
	}

	h.render(w, r, "admin_halal.tmpl", &templateData{
		User: h.currentUser(r),
		Data: items,
	})
}

func (h *Handler) adminHalalApprove(w http.ResponseWriter, r *http.Request) {
	h.adminHalalUpdate(w, r, "approved")
}

func (h *Handler) adminHalalReject(w http.ResponseWriter, r *http.Request) {
	h.adminHalalUpdate(w, r, "rejected")
}

func (h *Handler) adminHalalUpdate(w http.ResponseWriter, r *http.Request, status string) {
	if err := r.ParseForm(); err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	parts := splitPath(r.URL.Path)
	if len(parts) < 3 {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	oid, err := primitive.ObjectIDFromHex(parts[2])
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	note := strings.TrimSpace(r.PostForm.Get("note"))

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	v, err := h.App.Halal.FindByID(ctx, oid)
	if err != nil {
		h.clientError(w, http.StatusNotFound)
		return
	}

	if err := h.App.Halal.UpdateStatus(ctx, oid, status, note); err != nil {
		h.serverError(w, err)
		return
	}

	rest, err := h.App.Restaurants.FindByID(ctx, v.RestaurantID)
	if err == nil {
		if status == "approved" {
			rest.HalalStatus = "verified"
		}
		if status == "rejected" {
			rest.HalalStatus = "rejected"
		}
		_ = h.App.Restaurants.Update(ctx, rest)
	}

	http.Redirect(w, r, "/admin/halal", http.StatusSeeOther)
}

func (h *Handler) adminReviewsList(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	reviews, err := h.App.Reviews.ListAll(ctx)
	if err != nil {
		h.serverError(w, err)
		return
	}

	h.render(w, r, "admin_reviews.tmpl", &templateData{
		User: h.currentUser(r),
		Data: map[string]any{
			"reviews": reviews,
		},
	})
}

func (h *Handler) adminReviewHide(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	parts := splitPath(r.URL.Path)
	if len(parts) < 3 {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	oid, err := primitive.ObjectIDFromHex(parts[2])
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	note := strings.TrimSpace(r.PostForm.Get("note"))

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.App.Reviews.SetHidden(ctx, oid, true, u.ID, note); err != nil {
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/admin/reviews", http.StatusSeeOther)
}

func (h *Handler) adminReviewShow(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	parts := splitPath(r.URL.Path)
	if len(parts) < 3 {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	oid, err := primitive.ObjectIDFromHex(parts[2])
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	note := strings.TrimSpace(r.PostForm.Get("note"))

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.App.Reviews.SetHidden(ctx, oid, false, u.ID, note); err != nil {
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/admin/reviews", http.StatusSeeOther)
}
